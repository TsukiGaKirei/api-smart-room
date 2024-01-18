package schema

import "time"

type ResponseWebView struct {
	UserRoomWebView []UserRoomWebView `json:"user_room_data"`
	UserWebView     []UserWebView     `json:"user_data"`
	RoomsWebview    []RoomsWebview    `json:"rooms_data"`
	MqttLog         []MqttLog         `json:"mqtt_log"`
}

type RoomsWebview struct {
	Rid           int     `json:"rid"`
	Name          string  `json:"name"`
	Lamp          bool    `json:"lamp"`
	Ac            bool    `json:"ac"`
	AcTemp        int     `json:"ac_temp"`
	RoomTemp      float32 `json:"room_temp"`
	LastUpdated   string  `json:"last_updated"`
	LastUpdatedBy int     `json:"last_updated_by"`
	Door          bool    `json:"door"`
}

type UserRoomWebView struct {
	Uid              int     `json:"uid"`
	Rid              int     `json:"rid"`
	Distance         int     `json:"distance"`
	DesiredThreshold int     `json:"desired_threshold"`
	Threshold        float32 `json:"threshold"`
	LastUpdated      string  `json:"last_updated"`
}
type UserWebView struct {
	UID                 int     `json:"uid"`
	Name                string  `json:"name"`
	DesiredRadius       float32 `json:"desired_radius"`
	DesiredTemp         int     `json:"desired_temp"`
	DesiredThreshold    int     `json:"desired_threshold"`
	SmartRoomAutomation bool    `json:"smart_room_automation"`
	LastUpdated         string  `json:"last_updated"`
}

type MqttLog struct {
	Id          int       `json:"id"`
	Topic       string    `json:"topic"`
	Message     string    `json:"message"`
	PublishedAt time.Time `json:"published_at"`
}
