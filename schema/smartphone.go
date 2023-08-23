package schema

// Preference payload -> id_role, opendoor when inside radius, desired temp, deisred radius
type Preference struct {
	IdUser        int `gorm:"primaryKey"`
	DesiredRadius float32
	DesiredTemp   float32
	Autolockdoor  bool
}

type LocationUpdate struct {
	IdUser    string `json:"id_user"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}
