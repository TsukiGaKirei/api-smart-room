package schema

type Users struct {
	UID                 int     `json:"uid"`
	Name                string  `json:"name"`
	Longitude           float32 `json:"longitude"`
	Latitude            float32 `json:"latitude"`
	DesiredRadius       float32 `json:"desired_radius"`
	DesiredTemp         int     `json:"desired_temp"`
	SmartRoomAutomation bool    `json:"smart_room_automation"`
}
