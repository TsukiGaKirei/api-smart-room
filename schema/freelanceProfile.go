package schema

import "api-smart-room/model"

type FreelanceProfile struct {
	Nama      string               `json:"nama"`
	Email     string               `json:"email"`
	IdUserNik string               `json:"id_user_nik"`
	NlpTags   *model.FreelancerNlp `json:"nlp_tags"`
	Keahlian  string               `json:"keahlian"`
	Points    float64              `json:"points"`
	Alamat    string               `json:"alamat"`
}
