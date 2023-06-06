package schema

import "api-smart-room/model"

type OrderArrangement struct {
	ValueClean int64             `json:"harga"`
	Tasks      []model.OrderTask `json:"tasks"`
}
