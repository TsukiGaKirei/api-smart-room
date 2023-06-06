package model

import "time"

type Role struct {
	IdRole    string `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Role) TableName() string {
	return "public.role"
}
