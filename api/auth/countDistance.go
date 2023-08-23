package auth

import (
	"api-smart-room/model"
	"api-smart-room/schema"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// CountDistanceMapsApi function get distance then
func CountDistanceMapsApi(c echo.Context, o schema.Coordinates, d []schema.RoomCoordinates) error {
	// parameter can
	//var result map[schema.RoomCoordinates]float32
	var output model.DistanceMatrixResponse
	originsCoordinate := fmt.Sprintf("%f%s%f", o.Latitude, "%2C", o.Longitude)
	destinationCoordinate := fmt.Sprintf("%f%s%f", d[0].Latitude, "%2C", d[0].Longitude)
	if len(destinationCoordinate) >= 1 {
		for i := 1; i < len(d); i++ {
			addCoordinate := fmt.Sprintf("%s%d%s%d", "%7C", d[i].Latitude, "%2C", d[i].Longitude)
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}

	jsonErr := json.Unmarshal(body, &output)
	if jsonErr != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, output)
}

// UpdateMicroController Post location, room, desired temp, desired radius
//
func UpdateMicroController(c echo.Context) error {

	userID := c.QueryParam("id")
	if userID == "0" {
		userID = "hubla"

	}
	return c.JSON(http.StatusOK, userID)
}
