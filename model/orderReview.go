package model

import "time"

type OrderReview struct {
	IdReview    int       `gorm:"primaryKey;autoIncrement;" json:"id_review"`
	IdOrder     string    `json:"id_order"`
	IdFreelance int       `json:"id_freelance"`
	Rating      float64   `json:"rating"`
	NlpScore    float64   `json:"nlp_score"`
	PointReview float64   `json:"point_review"`
	Commentary  string    `json:"commentary"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
