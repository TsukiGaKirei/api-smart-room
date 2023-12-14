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

	"github.com/labstack/echo/v4"
)

// DistanceInfo represents the extracted information from the Distance Matrix API response
type DistanceInfo struct {
	RID      int `json:"rid"`
	Distance int `json:"distance"`
}

func CountDistanceMapsApi(c echo.Context, o schema.UserCoordinates) error {
	var apiResponse model.DistanceMatrixResponse
	var d []schema.RoomCoordinates
	db := database.GetDBInstance()
	if err := db.Raw(`
	select r.rid,r.latitude,r.longitude from rooms r, users_rooms ur
	where r.rid = ur.rid and ur.uid = ?
	`, o.UID).Scan(&d).Error; err == nil {
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
	// Extract distance information
	distanceInfo := make([]DistanceInfo, 0)

	for i, row := range d {
		var read DistanceInfo
		read.RID = row.RID
		read.Distance = apiResponse.Rows[0].Elements[i].Distance.Meters

		distanceInfo = append(distanceInfo, read)
	}

	//if user inside radius when  the room already active then threshold 0
	// if a user activated a room via automation then threshold 0 vice versa.
	// vice versa if user outside radius when the room is inactive then threshold 0.

	return c.JSON(http.StatusOK, apiResponse)
}
