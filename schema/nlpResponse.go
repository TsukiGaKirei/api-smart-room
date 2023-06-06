package schema

type ParentNlpApiResponse struct {
	Data NlpApiResponse `json:"data"`
}

type NlpApiResponse struct {
	NlpScore       float64 `json:"nlp_score"`
	RatingModelSum float64 `json:"rating_model_sum"`
}

type ParentNlpTag struct {
	Data []NlpApiResponse `json:"data"`
}

type NlpTagResp struct {
	Sifat string  `json:"sifat"`
	Value float64 `json:"value"`
}
