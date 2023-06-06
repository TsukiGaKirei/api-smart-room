package schema

type CoordinateBoth struct {
	FrLong float64 `json:"fr_long"`
	FrLat  float64 `json:"fr_lat"`
	ClLong float64 `json:"cl_long"`
	ClLat  float64 `json:"cl_lat"`
}
