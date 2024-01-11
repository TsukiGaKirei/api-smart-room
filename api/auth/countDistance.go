package auth

import (
	"api-smart-room/database"
	"api-smart-room/model"
	"api-smart-room/schema"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

// DistanceInfo represents the extracted information from the Distance Matrix API response

func CountDistanceMapsApi(c echo.Context, o schema.UserCoordinates) error {
	var apiResponse model.DistanceMatrixResponse
	var d []schema.RoomCoordinates
	db := database.GetDBInstance()
	if err := db.Raw(`
	select r.rid,r.latitude,r.longitude, ur.threshold, ur.last_updated, ur.distance, r.ac  , r.count_person, now() as current_time from rooms r, users_rooms ur, users u
	where r.rid = ur.rid and ur.uid = ? and u.uid=ur.uid and u.smart_room_automation=true
	`, o.UID).Scan(&d).Error; err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}

	originsCoordinate := fmt.Sprintf("%f%s%f", o.Latitude, "%2C", o.Longitude)
	destinationCoordinate := fmt.Sprintf("%f%s%f", d[0].Latitude, "%2C", d[0].Longitude)

	if len(d) >= 2 {
		for i := 1; i < len(d); i++ {
			addCoordinate := fmt.Sprintf("%s%f%s%f", "%7C", d[i].Latitude, "%2C", d[i].Longitude)
			destinationCoordinate += addCoordinate
		}
	}
	url := fmt.Sprintf(`https://maps.googleapis.com/maps/api/distancematrix/json?origins=%s&destinations=%s&key=%s`, originsCoordinate, destinationCoordinate, os.Getenv("API_KEY"))

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}

	jsonErr := json.Unmarshal(body, &apiResponse)
	if jsonErr != nil {
		return echo.ErrInternalServerError
	}
	// Extract distance 73
	fmt.Println(url)
	for i, row := range d {
		row.Distance = apiResponse.Rows[0].Elements[i].Distance.Meters
		if row.AC && row.Distance > o.DesiredRadius && row.CountPerson == 0 {
			row.Threshold += float32(row.CurrentTime.Sub(*row.LastUpdated).Minutes())

			if row.Threshold >= float32(o.Threshold) {
				row.AC = false
				row.Threshold = 0

				PublishMessage(strconv.Itoa(row.Rid) + ";ac_off")
				if err := db.Exec("update rooms set ac = false where rid = ?; ", row.Rid).Error; err != nil {
					fmt.Println(err)
					return echo.ErrInternalServerError
				}
			}

		} else if !row.AC && row.Distance <= o.DesiredRadius {

			row.Threshold += float32(row.CurrentTime.Sub(*row.LastUpdated).Minutes())
			if row.Threshold >= float32(o.Threshold) {
				row.AC = true
				row.Threshold = 0
				PublishMessage(strconv.Itoa(row.Rid) + ";ac_on;" + strconv.Itoa(o.DesiredTemp))
				if err := db.Exec("update rooms set ac = true, ac_temp = ? where rid = ?; ", o.DesiredTemp, row.Rid).Error; err != nil {
					fmt.Println(err)
					return echo.ErrInternalServerError
				}

			}
		}
		if row.AC && row.Distance <= o.DesiredRadius && row.Threshold != 0 || !row.AC && row.Distance > o.DesiredRadius && row.Threshold != 0 {
			row.Threshold = 0

		}
		if err := db.Exec("update users_rooms set threshold = ?, distance = ? ,last_updated=now() where uid = ? and rid = ?; ", row.Threshold, row.Distance, o.UID, row.Rid).Error; err != nil {
			fmt.Println(err)
			return echo.ErrInternalServerError
		}
	}

	//if user inside radius when  the room already active then threshold 0
	// if a user activated a room via automation then threshold 0 vice versa.
	// vice versa if user outside radius when the room is inactive then threshold 0.

	return c.JSON(http.StatusOK, d)
}
