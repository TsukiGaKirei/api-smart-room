package model

import "time"

type OrderPayment struct {
	IdPayment  int `gorm:"primaryKey;autoIncrement;"`
	IdOrder    string
	ValueClean int64
	AppFee     int64
	ValueTotal int64
	IdMethod   int
	IsPaid     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
