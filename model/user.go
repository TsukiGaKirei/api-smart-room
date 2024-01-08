package model

import (
	"time"

	"api-smart-room/database"
)

type User struct {
	IdUser    string `gorm:"primaryKey"`
	Name      string `validate:"required"`
	Email     string `validate:"required"`
	Password  string `validate:"required,gte=5,lte=30"`
	NoWa      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RoleScan struct {
	IdRole string
	IdUser string
}

func (User) TableName() string {
	return "public.user"
}

func (u *User) FindRoles() []RoleScan {
	var roles []RoleScan
	db := database.GetDBInstance()

	db.Model(&User{}).Select("public.user.id_user, roles.id_role").Where("public.user.id_user = ?", u.IdUser).
		Joins("left join user_role on user_role.id_user = user.id_user").
		Joins("left join roles on roles.id_role = user_role.id_role").
		Scan(&roles)

	return roles
}

func (u *User) HaveRole(roleId string) bool {
	var x RoleScan
	db := database.GetDBInstance()
	db.Model(&User{}).Select("public.user.id_user, user_role.id_role").
		Joins("left join user_role on user_role.id_user = public.user.id_user").
		Where("public.user.id_user = ?", u.IdUser).
		Where("user_role.id_role = ?", roleId).
		Scan(&x)

	return x.IdRole != ""
}
