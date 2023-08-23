package auth

import (
	"api-smart-room/api/auth"

	"github.com/labstack/echo/v4"
)

// api yang dibutuhkan
// microcontroller
// get Room status(need api for it) then implement it(code in arduino)
// post room_temp, post lamp_stats
//
// Smartphone
// POST Location periodically (when location is inside the radius and timer threshold reach 0 then activate room, else if outside location when room is active will automatically
// deactive the room when the timer threshold reach 0)
//, Post desired temp, Post desired radius for automatic activation, post desired timer to activate room.
func AuthSubRoute(group *echo.Group) {
	// smartphone
	group.PUT("/updateLocation", auth.UpdateLocation)

	// microcontroller

}
