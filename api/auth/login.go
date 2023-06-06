package auth

import (
	"api-smart-room/config"
	"api-smart-room/helper"
	"api-smart-room/schema"
	"api-smart-room/static"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	UId     string `json:"uid"`
	RoleId1 string `json:"role_id_1"`
	jwt.StandardClaims
}

func Login(c echo.Context) error {
	form := new(schema.Login)

	if err := c.Bind(form); err != nil {
		return err
	}

	if err := c.Validate(form); err != nil {
		return err
	}

	obj, errEmail := helper.FindByEmail(form.Email)
	rlFound, errRole := helper.FindRoleByName(form.Role)

	switch {
	case errEmail != nil, helper.IsPwdFalse(obj.Password, form.Password),
		errRole != nil, obj.HaveRole(rlFound.IdRole) == false:
		return &static.AuthError{}
	}

	claims := &JwtCustomClaims{
		obj.IdUser,
		rlFound.IdRole,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 8760).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, errJwt := token.SignedString(config.GetSignatureKey())

	if errJwt != nil {
		resError := static.ResponseError{
			Error:   true,
			Message: errJwt.Error(),
		}
		return c.JSON(http.StatusInternalServerError, resError)
	}

	rsp := static.ResponseToken{
		Error: false,
		Token: t,
	}

	return c.JSON(http.StatusOK, rsp)
}
