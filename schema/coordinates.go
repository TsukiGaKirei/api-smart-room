package schema

type Coordinates struct {
	Latitude  float32 `json:"Latitude"`
	Longitude float32 `json:"longitude"`
}
type RoomCoordinates struct {
	RID       int     `json:"rid"`
	Latitude  float32 `json:"Latitude"`
	Longitude float32 `json:"longitude"`
}
