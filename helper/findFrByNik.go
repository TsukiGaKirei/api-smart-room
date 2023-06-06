package helper

import (
	"api-smart-room/database"
	"api-smart-room/model"

	"gorm.io/gorm"
)

func IsFreelancerExist(nik string) (*model.FreelanceData, error) {
	db := database.GetDBInstance()
	var fr model.FreelanceData

	err := db.First(&fr, "nik = ?", nik).Error
	if err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	return &fr, nil
}
