package schema

type OfferingItem struct {
	IdOrderFr  string `json:"id_order"`
	IdStatus   int    `json:"id_status"`
	JobTitle   string `json:"job_title"`
	ClientName string `json:"client_name"`
	At         string `json:"at"`
	Status     string `json:"status"`
}
