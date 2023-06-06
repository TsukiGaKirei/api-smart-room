package client

import (
	"api-smart-room/static"
	"encoding/json"
	"net/http"

	"api-smart-room/database"

	"github.com/labstack/echo/v4"
)

func PaymentMethod(c echo.Context) error {

	type Result struct {
		Id_method    int    `json:"id_method"`
		Payment_name string `json:"payment_name"`
	}
	var result []Result
	db := database.GetDBInstance()
	err := db.Raw("select * from payment_method").Scan(&result).Error
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, result)
}

func OrderPayment(c echo.Context) error {
	//post
	type Payload struct {
		Id_method int    `json:"id_method"`
		Id_order  string `json:"id_order"`
	}
	var payload Payload
	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
		return echo.ErrBadRequest
	}
	db := database.GetDBInstance()
	err := db.Raw("update order_payment set is_paid=true , id_method = ? where id_order=?", payload.Id_method, payload.Id_order).Scan(&payload).Error
	if err != nil {
		return echo.ErrInternalServerError
	}
	msg := static.ResponseCreate{
		Error:   false,
		Message: "Tagihan Terbayar",
	}

	return c.JSON(http.StatusOK, msg)
}
