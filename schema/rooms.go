package schema

import "time"

type RoomUpdate struct {
	RoomId      int     `json:"room_id"`
	Temperature float32 `json:"temperature"`
	PersonCount int     `json:"person_count"`
}

type Rooms struct {
	RID         int       `json:"rid"`
	Name        string    `json:"name"`
	Longitude   float32   `json:"longitude"`
	Latitude    float32   `json:"latitude"`
	Lamp        bool      `json:"lamp"`
	Ac          bool      `json:"ac"`
	AcTemp      int       `json:"ac_temp"`
	RoomTemp    float32   `json:"room_temp"`
	LastUpdated time.Time `json:"last_updated"`
}

type UserRoom struct {
	UID int `json:"uid"`
	RID int `json:"rid"`
}
