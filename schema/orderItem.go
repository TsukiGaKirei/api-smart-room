package schema

import "time"

type OrderItem struct {
	JobChildName string    `json:"keahlian"`
	Name         string    `json:"nama_freelancer"`
	CreatedAt    time.Time `json:"tanggal_order"`
}
