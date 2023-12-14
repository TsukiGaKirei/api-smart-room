package schema

type UserCoordinates struct {
	UID       int     `json:"uid"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
type RoomCoordinates struct {
	RID       int     `json:"rid"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
