package helper

import (
	"api-smart-room/database"
	"api-smart-room/model"
)

func FindRoleByName(name string) (*model.Role, error) {
	db := database.GetDBInstance()
	var role model.Role
	search := db.Where("name = ?", name).First(&role)
	if search.Error != nil {
		return nil, search.Error
	}

	return &role, nil
}
