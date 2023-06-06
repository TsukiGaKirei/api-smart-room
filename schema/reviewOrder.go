package schema

type ReviewOrder struct {
	Rating     float64 `json:"rating" validate:"required"`
	Commentary string  `json:"komentar" validate:"required"`
}
