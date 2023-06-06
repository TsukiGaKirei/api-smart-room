package freelance

import (
	"api-smart-room/helper"
	"api-smart-room/model"
	"api-smart-room/schema"
	"api-smart-room/static"

	"fmt"
	"net/http"

	"strconv"
	"time"

	"api-smart-room/database"

	"github.com/dustin/go-humanize"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func OfferingList(c echo.Context) error {
	uId, _ := helper.ExtractToken(c)
	user, err := helper.FindByUId(uId)
	if err != nil {
		return err
	}

	fr, err := user.FindFreelanceAcc()
	if err != nil {
		return err
	}

	db := database.GetDBInstance()

	var orders []schema.OfferingItem
	db.Model(&model.Order{}).Select(`public.order.id_order as id_order_fr, public.order.id_client, public.order.id_freelance, public.order.created_at as at, 
				public.job_child_code.job_child_name as job_title, public.client_data.id_user, public.user.name as client_name, public.order_status.id_status`).
		Where(`public.order.id_freelance = ?`, fr.IdFreelance).
		Joins(`left join public.client_data on public.client_data.id_client = public.order.id_client`).
		Joins(`left join public.freelance_data on public.freelance_data.id_freelance = public.order.id_freelance`).
		Joins(`left join public.user on public.user.id_user = public.client_data.id_user`).
		Joins(`left join public.job_child_code on public.job_child_code.job_child_code = public.order.job_child_code`).
		Joins(`left join public.order_status on public.order_status.id_status = public.order.id_status`).
		Where(`public.order.id_status IN ?`, []int{1, 2, 4, 6}). // Diterima, proses, assigned
		Scan(&orders)
	if orders == nil {
		orders = []schema.OfferingItem{}
	}

	result := &static.ResponseSuccess{
		Error: false,
		Data:  orders,
	}

	return c.JSON(http.StatusOK, result)
}

func OfferingDetail(c echo.Context) error {
	idOrder := c.Param("id_order")

	uId, _ := helper.ExtractToken(c)
	user, err := helper.FindByUId(uId)
	if err != nil {
		return err
	}

	fr, err := user.FindFreelanceAcc()
	if err != nil {
		return err
	}

	db := database.GetDBInstance()

	var order schema.OfferingDetail

	type ClientFreelanceLatLong struct {
		ClientLat     float64 `json:"client_lat"`
		ClientLong    float64 `json:"client_long"`
		FreelanceLat  float64 `json:"freelance_lat"`
		FreelanceLong float64 `json:"freelance_long"`
	}
	var LatLong ClientFreelanceLatLong
	errLatLong := db.Raw(`select o.job_lat client_lat, o.job_long client_long , fd.address_lat freelance_lat, fd.address_long freelance_long
	from "order" o, client_data cd , freelance_data fd 
	where  o.id_client =cd.id_client and o.id_freelance =fd.id_freelance and o.id_order =?`, idOrder).Scan(&LatLong).Error
	if errLatLong != nil {
		return echo.ErrInternalServerError
	}
	type Response struct {
		IdOrderFr  string  `json:"id_order"`
		JobTitle   string  `json:"job_title"`
		ClientName string  `json:"client_name"`
		Keluhan    string  `json:"keluhan"`
		NoWaClient string  `json:"no_wa_client"`
		IdStatus   int     `json:"id_status"`
		Status     string  `json:"status"`
		Biaya      string  `json:"biaya"`
		Komentar   string  `json:"komentar"`
		Rating     string  `json:"rating"`
		JobLong    float64 `json:"longitude"`
		JobLat     float64 `json:"latitude"`
		Distance   string  `json:"jarak"`
	}

	distanceMatrixResponse, _ := helper.CountDistance(LatLong.ClientLat, LatLong.ClientLong, LatLong.FreelanceLat, LatLong.FreelanceLong)

	res := db.Model(&model.Order{}).Select(`public.order.id_order as id_order_fr, public.order.id_client, public.order.id_freelance, 
				public.order.job_description as keluhan, public.user.no_wa as no_wa_client, public.order.job_long, public.order.job_lat,
				public.job_child_code.job_child_name as job_title, public.client_data.id_user, 
				public.user.name as client_name, public.order_status.status_name as status, public.order_status.id_status,
				public.order_payment.value_clean as biaya, public.order_review.commentary as komentar,
				public.order_review.rating as rating`).
		Where(`public.order.id_freelance = ?`, fr.IdFreelance).
		Where(`public.order.id_order = ?`, idOrder).
		Joins(`left join public.client_data on public.client_data.id_client = public.order.id_client`).
		Joins(`left join public.freelance_data on public.freelance_data.id_freelance = public.order.id_freelance`).
		Joins(`left join public.user on public.user.id_user = public.client_data.id_user`).
		Joins(`left join public.job_child_code on public.job_child_code.job_child_code = public.order.job_child_code`).
		Joins(`left join public.order_status on public.order_status.id_status = public.order.id_status`).
		Joins(`left join public.order_payment on public.order_payment.id_order = public.order.id_order`).
		Joins(`left join public.order_review on public.order_review.id_order = public.order.id_order`).
		Scan(&order)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	biayaInt, _ := strconv.Atoi(order.Biaya)
	order.Biaya = humanize.Comma(int64(biayaInt))
	var fullResponse Response
	fullResponse.IdOrderFr = order.IdOrderFr
	fullResponse.JobTitle = order.JobTitle
	fullResponse.ClientName = order.ClientName
	fullResponse.Keluhan = order.Keluhan
	fullResponse.NoWaClient = order.NoWaClient
	fullResponse.IdStatus = order.IdStatus
	fullResponse.Status = order.Status
	fullResponse.Biaya = order.Biaya
	fullResponse.Komentar = order.Komentar
	fullResponse.Rating = order.Rating
	fullResponse.JobLong = order.JobLong
	fullResponse.JobLat = order.JobLat
	fullResponse.Distance = distanceMatrixResponse.Rows[0].Elements[0].Distance.HumanReadable
	result := &static.ResponseSuccess{
		Error: false,
		Data:  fullResponse,
	}

	return c.JSON(http.StatusOK, result)
}

func GetCoordinateBoth(c echo.Context) error {
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

	uId, _ := helper.ExtractToken(c)
	user, err := helper.FindByUId(uId)
	if err != nil {
		return err
	}

	fr, err := user.FindFreelanceAcc()
	if err != nil {
		return err
	}

	data := &schema.CoordinateBoth{
		FrLong: fr.AddressLong,
		FrLat:  fr.AddressLat,
		ClLong: order.JobLong,
		ClLat:  order.JobLat,
	}
	result := &static.ResponseSuccess{
		Error: false,
		Data:  data,
	}

	return c.JSON(http.StatusOK, result)
}

func ConfirmOffering(c echo.Context) error {
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

	order.IdStatus = 2
	if err := db.Save(&order).Error; err != nil {
		return err
	}

	response := &static.ResponseCreate{
		Error:   false,
		Message: "Order berhasil dikonfirmasi oleh Freelancer",
	}

	return c.JSON(http.StatusOK, response)
}

func RejectOffering(c echo.Context) error {
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

	order.IdStatus = 3
	if err := db.Save(&order).Error; err != nil {
		return err
	}

	response := &static.ResponseCreate{
		Error:   false,
		Message: "Order telah ditolak oleh Freelancer",
	}

	return c.JSON(http.StatusOK, response)
}

func ArrangeOffering(c echo.Context) error {
	idOrder := c.Param("id_order")
	form := new(schema.ArrangeOrder)

	if err := c.Bind(form); err != nil {
		return err
	}

	if err := c.Validate(form); err != nil {
		return err
	}

	db := database.GetDBInstance()
	var order model.Order
	res := db.First(&order, "id_order = ?", idOrder)
	if err := res.Error; err != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	// TODO: buat check kalau udah ada paymentnya di update aja
	valueClean := int64(form.Value)
	appFee := int64(float64(valueClean) * (float64(2) / 100))
	fmt.Println(appFee)
	valueTotal := valueClean + appFee
	ordPayment := order.GetPayment()
	if ordPayment != nil {
		ordPayment.ValueClean = valueClean
		ordPayment.AppFee = appFee
		ordPayment.ValueTotal = valueTotal
		db.Save(ordPayment)
	} else {
		timeNow := time.Now()
		newOdPayment := model.OrderPayment{
			IdOrder:    order.IdOrder,
			ValueClean: int64(form.Value),
			AppFee:     appFee,
			ValueTotal: valueTotal,
			IdMethod:   1,
			IsPaid:     false,
			CreatedAt:  timeNow,
			UpdatedAt:  timeNow,
		}

		err1 := db.Create(&newOdPayment).Error
		if err1 != nil {
			return err1
		}
	}
	order.IdStatus = 4
	db.Save(&order)

	response := static.ResponseCreate{
		Error:   false,
		Message: "Biaya dan pekerjaan berhasil ditentukan",
	}

	return c.JSON(http.StatusCreated, response)
}

func AddTask(c echo.Context) error {
	idOrder := c.Param("id_order")
	form := new(schema.ArrangeTask)

	if err := c.Bind(form); err != nil {
		return err
	}

	if err := c.Validate(form); err != nil {
		return err
	}

	timeNow := time.Now()
	db := database.GetDBInstance()
	var order model.Order
	res := db.First(&order, "id_order = ?", idOrder)
	if err := res.Error; err != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	newOdTask := model.OrderTask{
		IdOrder:    order.IdOrder,
		TaskDesc:   form.Task,
		TaskStatus: false,
		CreatedAt:  timeNow,
		UpdatedAt:  timeNow,
	}
	_ = db.Transaction(func(tx *gorm.DB) error {
		tx.Create(&newOdTask)
		return nil
	})

	response := static.ResponseSuccess{
		Error: false,
		Data:  newOdTask,
	}

	return c.JSON(http.StatusCreated, response)
}

func DeleteTask(c echo.Context) error {
	idOrder := c.Param("id_order")
	idTask := c.Param("id_task")

	db := database.GetDBInstance()
	var order model.Order
	res := db.First(&order, "id_order = ?", idOrder)
	if err := res.Error; err != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	if err1 := db.Delete(&model.OrderTask{}, idTask).Error; err1 != nil {
		return err1
	}

	response := static.ResponseCreate{
		Error:   false,
		Message: "Pekerjaan berhasil dihapus",
	}
	return c.JSON(http.StatusOK, response)
}

func GetArrangement(c echo.Context) error {
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

	ordPayment := order.GetPayment()
	var valueClean int64
	if ordPayment == nil {
		valueClean = 0
	} else {
		valueClean = ordPayment.ValueClean
	}

	obj := schema.OrderArrangement{
		ValueClean: valueClean,
		Tasks:      order.GetTasks(),
	}
	response := static.ResponseSuccess{
		Error: false,
		Data:  obj,
	}

	return c.JSON(http.StatusOK, response)
}

func RefreshStatus(c echo.Context) error {
	idOrder := c.Param("id_order")

	db := database.GetDBInstance()

	var status schema.RefreshStatus
	res := db.Model(&model.Order{}).Select(`public.order.id_status as is, public.order_status.id_status, public.order_status.status_name`).
		Where(`public.order.id_order = ?`, idOrder).
		Joins(`left join public.order_status on public.order_status.id_status = public.order.id_status`).
		Scan(&status)

	if err := res.Error; err != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	response := static.ResponseSuccess{
		Error: false,
		Data:  status,
	}
	return c.JSON(http.StatusOK, response)
}

func HistoriOffering(c echo.Context) error {
	uId, _ := helper.ExtractToken(c)
	user, err := helper.FindByUId(uId)
	if err != nil {
		return err
	}

	fr, err := user.FindFreelanceAcc()
	if err != nil {
		return err
	}

	db := database.GetDBInstance()

	var orders []schema.OfferingItem
	db.Model(&model.Order{}).Select(`public.order.id_order as id_order_fr, public.order.id_client, public.order.id_freelance, public.order.created_at as at, 
				public.job_child_code.job_child_name as job_title, public.client_data.id_user, public.user.name as client_name`).
		Where(`public.order.id_freelance = ?`, fr.IdFreelance).
		Where(`public.order.id_status IN ?`, []int{3, 5, 7}). // Selesai dan ditolak
		Joins(`left join public.client_data on public.client_data.id_client = public.order.id_client`).
		Joins(`left join public.freelance_data on public.freelance_data.id_freelance = public.order.id_freelance`).
		Joins(`left join public.user on public.user.id_user = public.client_data.id_user`).
		Joins(`left join public.job_child_code on public.job_child_code.job_child_code = public.order.job_child_code`).
		Scan(&orders)
	if orders == nil {
		orders = []schema.OfferingItem{}
	}

	result := &static.ResponseSuccess{
		Error: false,
		Data:  orders,
	}

	return c.JSON(http.StatusOK, result)
}
