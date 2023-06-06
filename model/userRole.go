package model

type UserRole struct {
	Id     int `gorm:"primaryKey;autoIncrement;"`
	IdUser string
	IdRole string
}
