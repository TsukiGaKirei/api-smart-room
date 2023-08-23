package schema

type RoomUpdate struct {
	RoomId      int     `json:"room_id"`
	Temperature float32 `json:"temperature"`
	PersonCount int     `json:"person_count"`
}
