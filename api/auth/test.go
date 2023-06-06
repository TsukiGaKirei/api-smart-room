package auth

import (
	"api-smart-room/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func TestMapsApi(c echo.Context) error {

	url := `https://maps.googleapis.com/maps/api/distancematrix/json?origins=40.6655101,-73.89188969999998&destinations=40.659569,-73.933783&key=`

	var output model.DistanceMatrixResponse

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
