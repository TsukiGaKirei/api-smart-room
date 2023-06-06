package middleware

import (
	"api-smart-room/database"
	"api-smart-room/helper"
	"api-smart-room/model"
	"api-smart-room/static"

	"github.com/labstack/echo/v4"
)

func CheckRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uId, role1 := helper.ExtractToken(c)
			var scRole model.Role
			db := database.GetDBInstance()
			// Kalau role yg di jwt sama kaya yang di param CheckRole
			if db.First(&scRole, "id_role = ?", role1); scRole.Name == role {
				// Kalau user nggapunya role yang ada di jwt
				if u, err := helper.IsUserExist(uId); err != nil ||
					!u.HaveRole(role1) {
					return echo.ErrNotFound
				}

				return next(c)
			}

			return &static.Unauthorized{}
		}
	}
}
