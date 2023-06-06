package auth

import (
	"api-smart-room/api/auth"

	"github.com/labstack/echo/v4"
)

func AuthSubRoute(group *echo.Group) {
	group.POST("/sign-up/client", auth.SignUp)
	group.POST("/sign-up/freelancer", auth.SignUpFr)
	group.POST("/login", auth.Login)
	group.GET("/test", auth.TestMapsApi)

}
