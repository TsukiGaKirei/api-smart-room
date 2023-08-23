package schema

type Users struct {
	uid           int     `json:"uid"'`
	name          string  `json:"name"`
	longitude     float32 `json:"longitude"`
	latitude      float32 `json:"latitude"`
	desiredRadius float32 `json:"desired_radius"`
	desiredTemp   int     `json:"desired_temp"`
	autodoorlock  bool    `json:"autodoorlock"`
}
