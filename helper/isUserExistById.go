package helper

import (
	"api-smart-room/database"
	"api-smart-room/model"

	"gorm.io/gorm"
)

func IsUserExist(uID string) (*model.User, error) {
	db := database.GetDBInstance()
	var user model.User

	err := db.First(&user, "id_user = ?", uID).Error
	if err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	return &user, nil
}
