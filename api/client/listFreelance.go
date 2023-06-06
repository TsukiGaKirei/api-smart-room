package client

import (
	"api-smart-room/helper"
	"api-smart-room/model"
	"api-smart-room/static"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"api-smart-room/database"

	"github.com/labstack/echo/v4"
)

type Response struct {
	IdFreelance       int       `json:"id_freelance"`
	Name              string    `json:"name"`
	IsTrainee         bool      `json:"is_trainee"`
	Rating            float64   `json:"rating"`
	JobDone           int       `json:"job_done"`
	DateJoin          time.Time `json:"date_join"`
	Jenis_kelamin     string    `json:"jenis_kelamin"`
	DistanceHaversign float64   `json:"distance_haversign"`
	JobChildName      string    `json:"job_child_name"`
	Address           string    `json:"address"`
	AddressLat        float64   `json:"address_lat"`
	AddressLong       float64   `json:"address_long"`
	Distance          string    `json:"distance"`
	JobCode           string    `json:"job_code"`
	JobChildCode      string    `json:"job_child_code"`
	NlpTag1           string    `json:"nlp_tag1"`
	NlpTag2           string    `json:"nlp_tag2"`
	NlpTag3           string    `json:"nlp_tag3"`
	NlpTag4           string    `json:"nlp_tag4"`
	NlpTag5           string    `json:"nlp_tag5"`
}

func ListFreelance(c echo.Context) error {

	job_code := c.Param("job_code")
	var result []Response

	type coordinate struct {
		AddressLat  float64 `json:"address_lat"`
		AddressLong float64 `json:"address_long"`
	}
	var clientLatLong coordinate
	userID, _ := helper.ExtractToken(c)
	db := database.GetDBInstance()
	errClient := db.Raw(`select address_lat, address_long from client_data where id_user = ?`, userID).Scan(&clientLatLong).Error
	if errClient != nil {
		return echo.ErrInternalServerError
	}
	// sort_by 1 / 2/ 3 /4
	// 1 == by
	err := db.Raw(`SELECT fd.id_freelance, u."name",fd.is_trainee,fd.rating,fd.job_done,fd.job_child_code,jc.job_code, case when fd.is_male = true then 'Pria' else 'Wanita' end as jenis_kelamin,
	 (6371 * acos( cos( radians(fd.address_lat) ) * cos( radians( ? ) ) *cos( radians( ? ) - radians(fd.address_long) ) 
	+ sin( radians(fd.address_lat) ) * sin( radians( ? ) )) ) as distance_haversign,
	jcc.job_child_name,fd.address,fd.address_lat,fd.address_long , (select fn.nlp_tag1 where fn.id_freelance =fd.id_freelance 
        ),(
        select fn.nlp_tag2
        where fn.id_freelance =fd.id_freelance 
        ), (
        select fn.nlp_tag3
        where fn.id_freelance =fd.id_freelance 
        ),(
        select fn.nlp_tag4
        where fn.id_freelance =fd.id_freelance 
        ),(
        select fn.nlp_tag5
        where fn.id_freelance =fd.id_freelance 
        )
	from freelance_data fd, job_child_code jcc ,job_code jc  , "user" u , freelancer_nlp fn
	where jcc.job_code  = jc.job_code and fd.job_child_code =jcc.job_child_code and u.id_user = fd.id_user and jc.job_code=? and
	(6371 * acos( cos( radians(fd.address_lat) ) * cos( radians( ? ) ) *cos( radians( ? ) - radians(fd.address_long) ) 
	+ sin( radians(fd.address_lat) ) * sin( radians( ? ) )) ) <10 
	order by distance_haversign asc, fd.rating desc`, clientLatLong.AddressLat, clientLatLong.AddressLong, clientLatLong.AddressLat, job_code, clientLatLong.AddressLat, clientLatLong.AddressLong, clientLatLong.AddressLat).Scan(&result)
	if err.Error != nil {
		return echo.ErrInternalServerError
	}

	url := `https://maps.googleapis.com/maps/api/distancematrix/json?origins=` + fmt.Sprintf("%f", clientLatLong.AddressLat) + `,` + fmt.Sprintf("%f", clientLatLong.AddressLong) + `&destinations=`
	api_key := "&key=" + os.Getenv("API_KEY")

	for i, data := range result {
		if i == 0 {
			url = url + fmt.Sprintf("%f", data.AddressLat) + `,` + fmt.Sprintf("%f", data.AddressLong)
		} else {
			url = url + `%7C` + fmt.Sprintf("%f", data.AddressLat) + `,` + fmt.Sprintf("%f", data.AddressLong)
		}
	}
	url += api_key

	// // akses gmaps distance api
	var output model.DistanceMatrixResponse
	fmt.Println(url)
	client := &http.Client{}
	req, erraaaa := http.NewRequest("GET", url, nil)
	if erraaaa != nil {
		fmt.Println(erraaaa)
	}
	resMaps, erraaaa := client.Do(req)
	if erraaaa != nil {
		fmt.Println(erraaaa)
	}
	defer resMaps.Body.Close()

	body, erraaaa := ioutil.ReadAll(resMaps.Body)
	if erraaaa != nil {
		fmt.Println(erraaaa)
		return erraaaa
	}

	jsonErr := json.Unmarshal(body, &output)
	if jsonErr != nil {
		return jsonErr
	}
	//

	var resultJson []Response
	resultJson = make([]Response, 0)
	for _, data := range output.Rows {
		for i, outputData := range data.Elements {
			result[i].Distance = outputData.Distance.HumanReadable
			if outputData.Distance.Meters <= 10000 {
				resultJson = append(resultJson, result[i])
			}
		}
	}

	res := static.ResponseSuccess{
		Error: false,
		Data:  resultJson,
	}
	return c.JSON(http.StatusOK, res)
}
