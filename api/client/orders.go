package client

import (
	"api-smart-room/helper"
	"api-smart-room/model"
	"api-smart-room/schema"
	"api-smart-room/static"
	"encoding/json"
	"net/http"
	"time"

	"api-smart-room/database"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SubmitOrder(c echo.Context) error {
	form := new(schema.OrderSubmit)
	idFreelance := c.Param("id_freelance")

	if err := c.Bind(form); err != nil {
		return err
	}

	if err := c.Validate(form); err != nil {
		return err
	}

	db := database.GetDBInstance()
	uId, _ := helper.ExtractToken(c)
	user, errUid := helper.FindByUId(uId)
	if errUid != nil {
		return errUid
	}
	timeNow := time.Now()

	idOd := "OD-" + helper.RandomStr(8)
	if idOd == "" {
		return echo.ErrNotFound
	}

	clientData, errClient := user.FindClientAcc()
	if errClient != nil {
		return gorm.ErrRecordNotFound
	}

	var freelanceData model.FreelanceData
	if err := db.Where("id_freelance = ?", idFreelance).First(&freelanceData).Error; err != nil {
		return err
	}

	// Nanti disini bakalan ditambahkan api cari address (text)
	// yang didapatkan dari api Google Map
	// Param: longitude latitude, Response: Alamat

	type ClientLatLong struct {
		Latitude  float64 `json:"address_lat"`
		Longitude float64 `json:"address_long"`
	}

	var clientCoordinate ClientLatLong
	errClientCoordinate := db.Raw(`select  cd.address_lat ,cd.address_long  from client_data cd where cd.id_user =?`, uId).Scan(&clientCoordinate).Error
	if errClientCoordinate != nil {
		return echo.ErrInternalServerError
	}

	errOrder := db.Create(&model.Order{
		IdOrder:        idOd,
		IdClient:       clientData.IdClient,
		IdFreelance:    freelanceData.IdFreelance,
		JobChildCode:   freelanceData.JobChildCode,
		JobLong:        clientCoordinate.Longitude,
		JobLat:         clientCoordinate.Latitude,
		JobDescription: form.JobDescription,
		AlreadyPaid:    false,
		IdStatus:       1,
		CreatedAt:      timeNow,
		UpdatedAt:      timeNow,
	}).Error

	if errOrder != nil {
		return errOrder
	}

	res := static.ResponseSuccess{
		Error: false,
		Data:  idOd,
	}
	return c.JSON(http.StatusCreated, res)
}

func DetailPesanan(c echo.Context) error {
	idOrder := c.Param("id_order")

	db := database.GetDBInstance()

	var order schema.OrderDetail
	res := db.Model(&model.Order{}).Select(`public.order.job_description, public.freelance_data.rating,
			public.user.name, public.user.no_wa, public.job_child_code.job_child_name, public.order_status.status_name,
			public.order_payment.value_clean, public.order_payment.value_total, public.order.id_status`).
		Where(`public.order.id_order = ?`, idOrder).
		Joins(`left join public.freelance_data on public.freelance_data.id_freelance = public.order.id_freelance`).
		Joins(`left join public.user on public.user.id_user = public.freelance_data.id_user`).
		Joins(`left join public.job_child_code on public.job_child_code.job_child_code = public.order.job_child_code`).
		Joins(`left join public.order_status on public.order_status.id_status = public.order.id_status`).
		Joins(`left join public.order_payment on public.order_payment.id_order = public.order.id_order`).
		Scan(&order)

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if res.Error != nil {
		return res.Error
	}

	return c.JSON(http.StatusOK, order)
}

func ConfirmOrder(c echo.Context) error {
	idOrder := c.Param("id_order")

	db := database.GetDBInstance()
	var order model.Order
	res := db.First(&order, "id_order = ?", idOrder)
	if err := res.Error; err != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	timeNow := time.Now()
	order.IdStatus = 6
	order.StartAt = timeNow
	db.Save(&order)

	response := static.ResponseCreate{
		Error:   false,
		Message: "Order berhasil dikonfirmasi. Order akan segera diproses",
	}
	return c.JSON(http.StatusOK, response)
}

func CancelOrder(c echo.Context) error {
	idOrder := c.Param("id_order")

	db := database.GetDBInstance()
	var order model.Order
	res := db.First(&order, "id_order = ?", idOrder)
	if err := res.Error; err != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	order.IdStatus = 5
	db.Save(&order)

	response := static.ResponseCreate{
		Error:   false,
		Message: "Order telah dibatalkan oleh Client",
	}
	return c.JSON(http.StatusOK, response)
}

func FinishOrder(c echo.Context) error {
	idOrder := c.Param("id_order")

	db := database.GetDBInstance()
	var order model.Order
	res := db.First(&order, "id_order = ?", idOrder)
	if err := res.Error; err != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	timeNow := time.Now()
	order.IdStatus = 7
	order.FinishedAt = timeNow
	db.Save(&order)

	response := static.ResponseCreate{
		Error:   false,
		Message: "Order telah diselesaikan oleh Client",
	}
	return c.JSON(http.StatusOK, response)
}

func TasksList(c echo.Context) error {
	idOrder := c.Param("id_order")

	db := database.GetDBInstance()
	var order model.Order
	res := db.First(&order, "id_order = ?", idOrder)
	if err := res.Error; err != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	var frData model.FreelanceData
	if err := db.Find(&frData, "id_freelance = ?", order.IdFreelance).Error; err != nil {
		return err
	}
	frData.JobDone++
	db.Save(&frData)

	response := static.ResponseSuccess{
		Error: false,
		Data:  order.GetTasks(),
	}

	return c.JSON(http.StatusOK, response)
}

func HistoryOrder(c echo.Context) error {
	uId, _ := helper.ExtractToken(c)
	user, err := helper.FindByUId(uId)
	if err != nil {
		return err
	}

	cl, err := user.FindClientAcc()
	if err != nil {
		return err
	}

	db := database.GetDBInstance()
	var orders []schema.OrderItem
	db.Model(&model.Order{}).Select(`public.order.created_at, public.job_child_code.job_child_name,
			public.user.name, public.order_status.status_name`).
		Where(`public.order.id_client = ?`, cl.IdClient).
		Where(`public.order.id_status IN ?`, []int{3, 5, 7}).
		Joins(`left join public.job_child_code on public.job_child_code.job_child_code = public.order.job_child_code`).
		Joins(`left join public.client_data on public.client_data.id_client = public.order.id_client`).
		Joins(`left join public.user on public.user.id_user = public.client_data.id_user`).
		Joins(`left join public.order_status on public.order_status.id_status = public.order.id_status`).
		Scan(&orders)

	if orders == nil {
		orders = []schema.OrderItem{}
	}

	result := &static.ResponseSuccess{
		Error: false,
		Data:  orders,
	}

	return c.JSON(http.StatusOK, result)
}

// func ReviewOrder(c echo.Context) error {
// 	type review_order struct {
// 		Id_order   string `json:"id_order"`
// 		Rating     int    `json:"rating"`
// 		Commentary string `json:"commentary"`
// 	}
// 	var payload review_order
// 	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
// 		return echo.ErrBadRequest
// 	}
// 	db := database.GetDBInstance()
// 	var id_freelance string

// 	err := db.Raw(`select id_freelance from "order" where id_order = ?`, payload.Id_order).Scan((&id_freelance)).Error
// 	if err != nil {
// 		return echo.ErrInternalServerError
// 	}

// 	err2 := db.Raw(`insert into order_review(id_order,id_freelance,rating,commentary,created_at,updated_at) values(?,?,?,?,?,?)`,
// 		payload.Id_order, id_freelance, payload.Rating, payload.Commentary, time.Now(), time.Now())
// 	if err2.Error != nil {
// 		return echo.ErrInternalServerError
// 	}

// 	msg := static.ResponseCreate{
// 		Error:   false,
// 		Message: "Review Berhasil",
// 	}

// 	return c.JSON(http.StatusOK, msg)
// }

func ReportViolation(c echo.Context) error {
	type report_violation struct {
		Id_order string `json:"id_order"`
		Title    string `json:"title"`
		Desc     string `json:"desc"`
	}
	var payload report_violation
	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
		return echo.ErrBadRequest
	}
	db := database.GetDBInstance()

	userId, _ := helper.ExtractToken(c)

	err := db.Raw(`insert into report_violation(id_order,created_by,title,desc,report_status,created_at,updated_at)
	values(?,?,?,?,1,?,?)`, payload.Id_order, userId, payload.Title, payload.Desc, time.Now(), time.Now()).Error
	if err != nil {
		return echo.ErrInternalServerError
	}

	msg := static.ResponseCreate{
		Error:   false,
		Message: "Laporan Sukses",
	}

	return c.JSON(http.StatusCreated, msg)
}
