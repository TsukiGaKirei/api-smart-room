package schema

type OrderSubmit struct {
	JobDescription string `json:"job_description" validate:"required"`
}
