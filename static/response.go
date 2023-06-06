package static

type ResponseSuccess struct {
	Error bool        `json:"error"`
	Data  interface{} `json:"data"`
}

type ResponseCreate struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type ResponseError struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type ResponseToken struct {
	Error bool   `json:"error"`
	Token string `json:"token"`
}
