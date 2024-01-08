package schema

import "time"

type UserCoordinates struct {
	UID           int     `json:"uid"`
	Latitude      float32 `json:"latitude"`
	Longitude     float32 `json:"longitude"`
	Threshold     int     `json:"threshold"`
	DesiredRadius int     `json:"desired_radius"`
	DesiredTemp   int     `json:"desired_temp"`
}
type RoomCoordinates struct {
	Rid         int        `json:"rid"`
	Latitude    float32    `json:"latitude"`
	Longitude   float32    `json:"longitude"`
	Threshold   float32    `json:"threshold"`
	Distance    int        `json:"distance"`
	LastUpdated *time.Time `json:"last_updated"`
	AC          bool       `json:"ac"`
	CurrentTime *time.Time `json:"current_time"`
	CountPerson int        `json:"count_person"`
}
