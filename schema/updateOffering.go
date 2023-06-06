package schema

type UpdateOffering struct {
	Status int `json:"status" validate:"required"`
}
