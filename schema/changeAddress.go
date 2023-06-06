package schema

type ChangeAddress struct {
	Address     string  `json:"address"`
	AddressLong float64 `json:"address_long"`
	AddressLat  float64 `json:"address_lat"`
}
