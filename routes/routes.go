package routes

import (
	authRoute "api-smart-room/routes/auth"
	// clientRoute "api-smart-room/routes/client"
	// freelanceRoute "api-smart-room/routes/freelance"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const API_ROUTE = "/api/"
const CSRFTokenHeader = "X-CSRF-Token"
const CSRFKey = "csrf"

func Init() *echo.Echo {
	e := echo.New()

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.POST, echo.DELETE},
	}))

	// Group Auth
	authGroup := e.Group(API_ROUTE + "auth")
	authRoute.AuthSubRoute(authGroup)

	// Group Client
	// clientGroup := e.Group(API_ROUTE+"client", customMiddleware.IsAuthenticated)
	// clientGroup.Use(customMiddleware.CheckRole("client"))
	// clientRoute.ClientSubRoute(clientGroup)

	// Freelance Group
	// freelanceGroup := e.Group(API_ROUTE+"freelancer", customMiddleware.IsAuthenticated)
	// freelanceGroup.Use(customMiddleware.CheckRole("freelancer"))
	// freelanceRoute.FreelanceSubRoute(freelanceGroup)

	return e
}
