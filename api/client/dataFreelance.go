package client

import (
	"api-smart-room/helper"
	"api-smart-room/model"
	"api-smart-room/schema"
	"api-smart-room/static"
	"encoding/json"
	"errors"
	"fmt"

	"io/ioutil"
	"net/http"
	"os"

	"api-smart-room/database"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func DataFreelance(c echo.Context) error {
	idFreelance := c.Param("id_freelance")
	db := database.GetDBInstance()

	var freelanceData model.FreelanceData
	err := db.First(&freelanceData, "id_freelance = ?", idFreelance).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	var user model.User
	err = db.First(&user, "id_user = ?", freelanceData.IdUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	bidang, errBidang := freelanceData.FindFreelanceBidang()
	if errBidang != nil {
		return errBidang
	}
	keahlian, errKeahlian := freelanceData.FindFreelanceKeahlian()
	if errKeahlian != nil {
		return errKeahlian
	}
	nlpTag, errNlp := freelanceData.FindNlpTag()
	if errNlp != nil {
		return errNlp
	}

	uId, _ := helper.ExtractToken(c)
	type clientCoordinate struct {
		Longitude float64 `json:"address_long"`
		Latitude  float64 `json:"address_lat"`
	}
	var clientLongLat clientCoordinate

	errClient := db.Raw(`select address_long, address_lat from client_data where id_user=?`, uId).Scan(&clientLongLat).Error
	if errClient != nil {
		return echo.ErrInternalServerError
	}

	// TEST

	url := `https://maps.googleapis.com/maps/api/distancematrix/json?origins=` + fmt.Sprintf("%f", clientLongLat.Latitude) + `,` + fmt.Sprintf("%f", clientLongLat.Longitude) + `&destinations=` + fmt.Sprintf("%f", freelanceData.AddressLat) + `,` + fmt.Sprintf("%f", freelanceData.AddressLong) + `&key=` + os.Getenv("API_KEY")

	var output model.DistanceMatrixResponse
	fmt.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
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
	// TEST

	data := &schema.FreelanceData{
		Nama:     user.Name,
		Bidang:   bidang,
		Alamat:   freelanceData.Address,
		Keahlian: keahlian,
		NlpTag:   nlpTag,
		Distance: output.Rows[0].Elements[0].Distance.HumanReadable,
	}
	if freelanceData.IsMale {
		data.JenisKelamin = "Pria"
	} else {
		data.JenisKelamin = "Wanita"
	}

	result := &static.ResponseSuccess{
		Error: false,
		Data:  data,
	}

	return c.JSON(http.StatusOK, result)
}

// // // data freelancer by gemi
// func DataFreelancer(c echo.Context) error {
// 	idFreelance := c.Param("id_freelance")
// 	db := database.GetDBInstance()
// 	type queryResult struct {
// 		IdFreelance  int       `json:"id_freelance"`
// 		IdUser       string    `json:"id_user"`
// 		IsTrainee    bool      `json:"is_trainee"`
// 		Rating       float64   `json:"rating"`
// 		JobDone      int       `json:"job_done"`
// 		DateJoin     time.Time `json:"date_join"`
// 		Address      string    `json:"address"`
// 		AddressLong  float64   `json:"address_long"`
// 		AddressLat   float64   `json:"address_lat"`
// 		IsMale       bool      `json:"is_male"`
// 		Dob          time.Time `json:"dob"`
// 		Nik          string    `json:"nik"`
// 		ProfilePict  string    `json:"profile_pict"`
// 		Points       float64   `json:"points"`
// 		JobChildCode string    `json:"job_child_code"`
// 		CreatedAt    time.Time `json:"created_at"`
// 		UpdatedAt    time.Time `json:"updated_at"`
// 	}
// 	var result queryResult
// 	err := db.Raw(`select * from freelance_data,"user" u where fd.id_user = u.id_user and fd.id_freelance=?`, idFreelance).Scan(&result).Error
// 	if err != nil {
// 		return echo.ErrInternalServerError
// 	}

// 	res := &static.ResponseSuccess{
// 		Error: false,
// 		Data:  result,
// 	}

// 	return c.JSON(http.StatusOK, res)
// }
