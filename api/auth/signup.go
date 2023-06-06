package auth

import (
	"api-smart-room/database"
	"api-smart-room/helper"
	"api-smart-room/model"
	"api-smart-room/schema"
	"api-smart-room/static"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SignUp(c echo.Context) error {
	form := new(schema.SignUp)
	if err := c.Bind(form); err != nil {
		return err
	}

	if err := c.Validate(form); err != nil {
		return err
	}
	userExist, _ := helper.FindByEmail(form.Email)
	freelancerExist, _ := helper.IsFreelancerExist(form.Nik)
	if userExist != nil || freelancerExist != nil {
		msg := static.ResponseCreate{
			Error:   true,
			Message: "Email already exist",
		}
		return c.JSON(http.StatusBadRequest, msg)
	}

	inputLat, _ := strconv.ParseFloat(form.Latitude, 64)
	inputLong, _ := strconv.ParseFloat(form.Longitude, 64)

	db := database.GetDBInstance()

	timeNow := time.Now()
	newUser := &model.User{
		IdUser:    "CL" + "-" + helper.RandomStr(10),
		Name:      form.Nama,
		Email:     form.Email,
		NoWa:      form.NoWa,
		Password:  helper.GeneratePwd(form.Password),
		CreatedAt: timeNow, UpdatedAt: timeNow,
	}
	_ = db.Transaction(func(tx *gorm.DB) error {
		tx.Create(&newUser)
		return nil
	})

	valJenisKelamin := strings.ToLower(form.JenisKelamin)
	var isMale = false
	if valJenisKelamin == "pria" || valJenisKelamin == "laki-laki" || valJenisKelamin == "laki laki" || valJenisKelamin == "cowo" || valJenisKelamin == "cowok" || valJenisKelamin == "male" || valJenisKelamin == "jantan" {
		isMale = true
	} else if valJenisKelamin == "cewe" || valJenisKelamin == "cewek" || valJenisKelamin == "perempuan" || valJenisKelamin == "wanita" || valJenisKelamin == "betina" || valJenisKelamin == "female" {
		isMale = false
	}

	obj := db.Create(&model.ClientData{
		IdUser:      newUser.IdUser,
		Address:     form.Alamat,
		IsMale:      isMale,
		Nik:         form.Nik,
		AddressLong: inputLat,
		AddressLat:  inputLong,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	})
	clientRole, _ := helper.FindRoleByName("client")
	uR := db.Create(&model.UserRole{
		IdUser: newUser.IdUser,
		IdRole: clientRole.IdRole,
	})
	if obj.Error != nil || uR.Error != nil {
		return obj.Error
	}

	msg := static.ResponseCreate{
		Error:   false,
		Message: "Pengguna berhasil mendaftar",
	}

	return c.JSON(http.StatusCreated, msg)
}
