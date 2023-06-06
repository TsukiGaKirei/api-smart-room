package freelance

import (
	"api-smart-room/helper"
	"api-smart-room/schema"
	"api-smart-room/static"

	"net/http"

	"api-smart-room/database"

	"github.com/labstack/echo/v4"
)

func GetProfile(c echo.Context) error {
	uId, _ := helper.ExtractToken(c)
	user, err := helper.FindByUId(uId)
	if err != nil {
		return err
	}

	fr, errFr := user.FindFreelanceAcc()
	if errFr != nil {
		return errFr
	}

	nlpTags, errNlpTag := fr.FindNlpTag()
	if errNlpTag != nil {
		return errNlpTag
	}

	keahlian, errKeahlian := fr.FindFreelanceKeahlian()
	if errKeahlian != nil {
		return errKeahlian
	}

	data := schema.FreelanceProfile{
		Nama:      user.Name,
		Email:     user.Email,
		IdUserNik: user.IdUser + " / " + fr.Nik,
		NlpTags:   nlpTags,
		Points:    fr.Points,
		Keahlian:  keahlian,
		Alamat:    fr.Address,
	}

	res := static.ResponseSuccess{
		Error: false,
		Data:  data,
	}

	return c.JSON(http.StatusOK, res)
}

func UpdateAddress(c echo.Context) error {
	uId, _ := helper.ExtractToken(c)
	user, err := helper.FindByUId(uId)
	if err != nil {
		return err
	}

	fr, errFr := user.FindFreelanceAcc()
	if errFr != nil {
		return errFr
	}

	form := new(schema.ChangeAddress)

	if err := c.Bind(form); err != nil {
		return err
	}

	if err := c.Validate(form); err != nil {
		return err
	}

	db := database.GetDBInstance()
	fr.Address = form.Address
	fr.AddressLat = form.AddressLat
	fr.AddressLong = form.AddressLong
	if err = db.Save(&fr).Error; err != nil {
		return err
	}

	res := static.ResponseCreate{
		Error:   false,
		Message: "Alamat berhasil diganti",
	}
	return c.JSON(http.StatusOK, res)
}
