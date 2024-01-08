package schema

// Preference payload -> id_role, opendoor when inside radius, desired temp, deisred radius
type Preference struct {
	IdUser        int     `gorm:"primaryKey" json:"id_user"`
	DesiredRadius float32 `json:"desired_radius"`
	DesiredTemp   float32 `json:"desired_temp"`
	Autolockdoor  bool    `json:"autolockdoor"`
}

type LocationUpdate struct {
	UID       int     `json:"uid"`
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

type OpenDoor struct {
	UID    int `json:"uid"`
	Radius int `json:"radius"`
}

type UserConfig struct {
	UID                 int  `json:"uid"`
	DesiredTemp         int  `json:"desired_temp"`
	DesiredRadius       int  `json:"desired_radius"`
	SmartRoomAutomation bool `json:"smart_room_automation"`
	DesiredThreshold    int  `json:"desired_threshold"`
}
