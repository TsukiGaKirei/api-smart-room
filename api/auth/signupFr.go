package auth

import (
	"api-smart-room/helper"
	"api-smart-room/model"
	"api-smart-room/schema"
	"api-smart-room/static"
	"net/http"
	"strconv"
	"time"

	"api-smart-room/database"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SignUpFr(c echo.Context) error {
	form := new(schema.SignUpFreelance)
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
			Message: "Pengguna telah terdaftar sebelumnya.",
		}
		return c.JSON(http.StatusBadRequest, msg)
	}

	db := database.GetDBInstance()

	if err := db.First(&model.JobChildCode{}, "job_child_code = ?", form.JobChildCode).Error; err != nil {
		return err
	}

	timeNow := time.Now()
	newUser := &model.User{
		IdUser:    "FR" + "-" + helper.RandomStr(10),
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

	convertJk, _ := strconv.ParseBool(form.JenisKelamin)
	formatParse := "2012-12-30"
	dobParse, _ := time.Parse(formatParse, form.Dob)
	obj := db.Create(&model.FreelanceData{
		IdUser:       newUser.IdUser,
		IsTrainee:    false,
		Address:      form.Address,
		AddressLong:  form.AddressLong,
		AddressLat:   form.AddressLat,
		Rating:       0,
		Points:       0,
		JobDone:      0,
		DateJoin:     timeNow,
		IsMale:       convertJk,
		Dob:          dobParse,
		Nik:          form.Nik,
		JobChildCode: form.JobChildCode,
		CreatedAt:    timeNow,
		UpdatedAt:    timeNow,
	})
	clientRole, _ := helper.FindRoleByName("freelancer")
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
