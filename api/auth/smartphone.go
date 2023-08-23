package auth

import (
	"api-smart-room/database"
	"api-smart-room/schema"
	"api-smart-room/static"
	"encoding/json"
	"strconv"
	"time"

	"net/http"

	"github.com/labstack/echo/v4"
)

func UpdateLocation(c echo.Context) error {
	var payload schema.LocationUpdate
	postLong, _ := strconv.ParseFloat(payload.Longitude, 64)
	postLat, _ := strconv.ParseFloat(payload.Latitude, 64)
	postId, _ := strconv.ParseInt(payload.IdUser, 4, 4)
	const updateSql = `
	UPDATE users
	SET longitude=?, latitude=?
	WHERE id=?
`
	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
		return echo.ErrBadRequest
	}
	db := database.GetDBInstance()
	err := db.Raw(updateSql, postLong, postLat, time.Now(), postId)
	if err == nil {
		return echo.ErrInternalServerError
	}
	res := static.ResponseSuccess{
		Error: false,
		Data:  "Location Updated"}

	return c.JSON(http.StatusCreated, res)
}
