package auth

import (
	"api-smart-room/database"
	"api-smart-room/schema"
	"api-smart-room/static"
	"encoding/json"
	"fmt"
	"strconv"

	"net/http"

	"github.com/labstack/echo/v4"
)

func UpdateLocation(c echo.Context) error {
	var payload schema.LocationUpdate
	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
		return echo.ErrBadRequest
	}

	const updateSql = `
	UPDATE users
	SET longitude=?, latitude=?
	WHERE uid=?
`

	db := database.GetDBInstance()
	err := db.Exec(updateSql, payload.Longitude, payload.Latitude, payload.UID).Error
	if err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}

	var UserCoordinates schema.UserCoordinates

	err = db.Raw(`select u.uid , u.latitude ,u.longitude ,u.desired_radius ,u.threshold , u.desired_temp 
	from users u 
	where u.uid =?`, payload.UID).Scan(&UserCoordinates).Error
	if err != nil {
		fmt.Println(err)

		return echo.ErrInternalServerError
	}
	CountDistanceMapsApi(c, UserCoordinates)

	res := static.ResponseSuccess{
		Error: false,
		Data:  "Location Updated"}

	return c.JSON(http.StatusCreated, res)
}

// due to device incapability, door can only be opened for 10 seconds maximum
func OpenDoor(c echo.Context) error {
	var payload schema.OpenDoor
	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
		return echo.ErrBadRequest
	}
	if payload.Radius == 0 && payload.UID == 0 {
		return echo.ErrBadRequest
	}

	var userRoom []schema.UserRoom
	const getRoomInsideRadius = `
	select ur.uid,ur.rid,ur.distance,ur.last_updated from users_rooms ur where ur.uid =? and distance <=?
`

	db := database.GetDBInstance()
	if err := db.Raw(getRoomInsideRadius, payload.UID, payload.Radius).Scan(&userRoom).Error; err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}
	fmt.Println(userRoom)

	for _, room := range userRoom {

		if room.Distance <= payload.Radius {
			PublishMessage(strconv.Itoa(room.Rid) + ";door_open")
			fmt.Println("Message published room -> " + strconv.Itoa(room.Rid) + "open_door")
		}
	}
	res := static.ResponseSuccess{
		Error: false,
		Data:  "Door Opened"}

	return c.JSON(http.StatusCreated, res)
}

// due to device incapability, door can only be opened for 10 seconds maximum
func UpdateConfiguration(c echo.Context) error {
	var payload schema.UserConfig
	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
		return echo.ErrBadRequest
	}
	var userRoom []schema.UserRoom
	var rid []int
	const updateRoomConfiguration = `
	update users set threshold=?,desired_temp = ?, desired_radius = ?, smart_room_automation=? where uid = ?;
`

	const selectUserRoom = `	
select ur.rid  
from users_rooms ur , rooms r, users u 
where ur.uid =? and  ur.rid =r.rid and u.uid = ur.uid and ur.distance <= u.desired_radius  and r.ac =true
	`
	//get activated room where the user is registred and online
	db := database.GetDBInstance()
	if err := db.Raw(updateRoomConfiguration, payload.DesiredThreshold, payload.DesiredTemp, payload.DesiredRadius, payload.SmartRoomAutomation, payload.UID).Scan(&userRoom).Error; err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}
	if err := db.Raw(selectUserRoom, payload.UID).Scan(&rid).Error; err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}
	if err := db.Exec(`update rooms set ac_temp = ? where rid in(
		select ur.rid  
		from users_rooms ur , rooms r, users u 
		where ur.uid =? and  ur.rid =r.rid and u.uid = ur.uid and ur.distance <= u.desired_radius  and r.ac =true)`, payload.DesiredTemp, payload.UID).Error; err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}
	for _, room := range rid {
		PublishMessage(strconv.Itoa(room) + ";ac_on;" + strconv.Itoa(payload.DesiredTemp))
		fmt.Println("Message published room -> " + strconv.Itoa(room) + "ac_temp;" + strconv.Itoa(payload.DesiredTemp))
	}
	fmt.Println(rid)
	res := static.ResponseSuccess{
		Error: false,
		Data:  "Configuration Updated"}

	return c.JSON(http.StatusCreated, res)
}
