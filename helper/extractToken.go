package helper

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func ExtractToken(c echo.Context) (string, string) {
	userLogged := c.Get("user")
	token := userLogged.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	uID := claims["uid"].(string)
	role1 := claims["role_id_1"].(string)

	return uID, role1
}
