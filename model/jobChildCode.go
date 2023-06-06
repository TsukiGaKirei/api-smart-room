package model

import "time"

type JobChildCode struct {
	JobChildCode string `gorm:"primaryKey"`
	JobChildName string
	JobCode      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
