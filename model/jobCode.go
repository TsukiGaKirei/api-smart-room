package model

import "time"

type JobCode struct {
	JobCode     string `gorm:"primaryKey"`
	JobCategory string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
