package model

type FreelancerNlp struct {
	IdFreelanceNlp int    `gorm:"primaryKey;autoIncrement;" json:"id_freelance_nlp"`
	IdFreelance    int    `json:"id_freelance"`
	NlpTag1        string `json:"nlp_tag1"`
	NlpTag2        string `json:"nlp_tag2"`
	NlpTag3        string `json:"nlp_tag3"`
	NlpTag4        string `json:"nlp_tag4"`
	NlpTag5        string `json:"nlp_tag5"`
}

func (FreelancerNlp) TableName() string {
	return "public.freelancer_nlp"
}
