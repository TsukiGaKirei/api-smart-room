package schema

type OrderDetail struct {
	Name           string `json:"nama"`
	JobChildName   string `json:"keahlian"`
	NoWa           string `json:"no_whatsapp"`
	Rating         int    `json:"rating"`
	IdStatus       int    `json:"id_status"`
	StatusName     string `json:"status"`
	JobDescription string `json:"keluhan"`
	ValueClean     int64  `json:"perkiraan_harga"`
	ValueTotal     int64  `json:"total"`
}
