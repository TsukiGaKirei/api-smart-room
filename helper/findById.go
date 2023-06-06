package helper

import (
	"api-smart-room/database"
	"api-smart-room/model"
)

func FindByUId(uId string) (*model.User, error) {
	db := database.GetDBInstance()
	var user model.User
	search := db.Where("id_user = ?", uId).First(&user)
	if search.Error != nil {
		return nil, search.Error
	}

	return &user, nil
}
