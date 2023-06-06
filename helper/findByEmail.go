package helper

import (
	"api-smart-room/database"
	"api-smart-room/model"
)

func FindByEmail(email string) (*model.User, error) {
	db := database.GetDBInstance()
	var user model.User
	search := db.Where("email = ?", email).First(&user)
	if search.Error != nil {
		return nil, search.Error
	}

	return &user, nil
}
