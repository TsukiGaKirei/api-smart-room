package middleware

import (
	"api-smart-room/config"

	"github.com/labstack/echo/v4/middleware"
)

var IsAuthenticated = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: config.GetSignatureKey(),
})
