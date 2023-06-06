package model

import "time"

type ClientData struct {
	IdClient    int    `gorm:"primaryKey;autoIncrement;"`
	IdUser      string `json:"id_user"`
	Address     string
	AddressLong float64 `json:"address_long"`
	AddressLat  float64 `json:"address_lat"`
	IsMale      bool
	Nik         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
