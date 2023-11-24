package helper

import (
	"api-smart-room/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//receive message-> update
func CountDistance(originLat float64, originLong float64, destinationLat float64, destinationLong float64) (*model.DistanceMatrixResponse, error) {
	url := `https://maps.googleapis.com/maps/api/distancematrix/json?origins=` + fmt.Sprintf("%f", originLat) + `,` + fmt.Sprintf("%f", originLong) + `&destinations=` + fmt.Sprintf("%f", destinationLat) + `,` + fmt.Sprintf("%f", destinationLong) + `&key=` + os.Getenv("API_KEY")

	var output model.DistanceMatrixResponse
	fmt.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &output, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return &output, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return &output, err
	}

	jsonErr := json.Unmarshal(body, &output)
	if jsonErr != nil {
		return &output, err
	}

	return &output, nil

}
