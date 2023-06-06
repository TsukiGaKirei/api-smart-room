package client

import (
	"api-smart-room/helper"
	"api-smart-room/static"
	"net/http"

	"api-smart-room/database"

	"github.com/labstack/echo/v4"
)

func DataPersonal(c echo.Context) error {
	type Result struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		// HousePict   string    `json:"house_pict"`
		NoWa        string  `json:"no_wa"`
		Address     string  `json:"address"`
		AddressLong float64 `json:"address_long"`
		AddressLat  float64 `json:"address_lat"`
		IsMale      bool    `json:"is_male"`
		// Dob         time.Time `json:"dob"`
		Nik string `json:"nik"`
		// ProfilePict string    `json:"profile_pict"`
	}
	type Response struct {
		Name  string `json:"nama"`
		Email string `json:"email"`
		// HousePict   string    `json:"house_pict"`
		NoWa         string  `json:"no_wa"`
		Address      string  `json:"alamat"`
		AddressLong  float64 `json:"address_long"`
		AddressLat   float64 `json:"address_lat"`
		JenisKelamin string  `json:"jenis_kelamin"`
		// Dob         time.Time `json:"dob"`
		Nik string `json:"nik"`
		// ProfilePict string    `json:"profile_pict"`
	}
	var result Result
	var jsonResp Response
	uId, _ := helper.ExtractToken(c)

	db := database.GetDBInstance()
	err := db.Raw(`select u.email,u.no_wa,u."name"  ,cd.address ,cd.address_long,cd.address_lat ,cd.is_male,cd.nik 
	from "user" u, client_data cd
	where u.id_user = cd.id_user and u.id_user=?`, uId).Scan(&result).Error
	if err != nil {
		return echo.ErrInternalServerError
	}
	jsonResp.Address = result.Address
	jsonResp.AddressLat = result.AddressLat
	jsonResp.AddressLong = result.AddressLong
	jsonResp.Name = result.Name
	jsonResp.Email = result.Email
	jsonResp.NoWa = result.NoWa
	jsonResp.Nik = result.Nik
	if result.IsMale {
		jsonResp.JenisKelamin = "Pria"
	} else {
		jsonResp.JenisKelamin = "wanita"
	}

	// err := db.Raw(`select u.email,u.no_wa,u."name"  ,cd.address ,cd.address_long,cd.address_lat ,cd.is_male,cd.nik ,cd.dob  ,cd.profile_pict ,cd.house_pict
	// from "user" u, client_data cd
	// where u.id_user = cd.id_user and u.id_user=?`, uId).Scan(&result).Error
	// if err != nil {
	// 	return echo.ErrInternalServerError
	// }
	res := static.ResponseSuccess{
		Error: false,
		Data:  jsonResp,
	}

	return c.JSON(http.StatusOK, res)
}
