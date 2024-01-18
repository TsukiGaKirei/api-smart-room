package auth

import (
	"api-smart-room/database"
	"api-smart-room/schema"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetDataAllRoom(c echo.Context) error {
	var result schema.ResponseWebView
	db := database.GetDBInstance()
	if err := db.Raw(`select r.rid ,r."name"  ,r.ac_temp , r.last_updated , r.ac, r.lamp ,r.door , r.last_updated_by 
	from rooms r `).Scan(&result.RoomsWebview).Error; err != nil {
		return echo.ErrInternalServerError
	}

	if err := db.Raw(`select ur.uid ,ur.rid , ur.distance ,u.threshold as desired_threshold, ur.threshold , ur.last_updated 
	from users_rooms ur ,users u 
	where ur.uid = u.uid `).Scan(&result.UserRoomWebView).Error; err != nil {
		return echo.ErrInternalServerError
	}

	if err := db.Raw(`select u.uid ,u."name" ,u.desired_radius ,u.desired_temp , u.threshold as desired_threshold,u.smart_room_automation , u.last_updated
	from users u  `).Scan(&result.UserWebView).Error; err != nil {
		return echo.ErrInternalServerError
	}

	if err := db.Raw(`select ml.id ,ml.topic ,ml.message ,ml.published_at 
	from mqtt_log ml
	order by ml.published_at desc
	`).Scan(&result.MqttLog).Error; err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, result)
}
