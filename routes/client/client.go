package client

import (
	"api-smart-room/api/client"

	"github.com/labstack/echo/v4"
)

func ClientSubRoute(group *echo.Group) {
	group.GET("/me", client.DataPersonal)
	group.GET("/search/freelance/:job_code", client.ListFreelance)

	group.GET("/freelance/:id_freelance", client.DataFreelance)
	group.GET("/payment/method", client.PaymentMethod)

	group.POST("/freelance/:id_freelance/order", client.SubmitOrder)
	group.GET("/orders/:id_order", client.DetailPesanan)
	group.PATCH("/orders/:id_order/confirm", client.ConfirmOrder)
	group.PATCH("/orders/:id_order/cancel", client.CancelOrder)
	group.PATCH("/orders/:id_order/finish", client.FinishOrder)
	group.PUT("/orders/payment", client.OrderPayment)
	group.POST("/orders/report", client.ReportViolation)

	group.GET("/orders/:id_order/tasks", client.TasksList)

	group.GET("/history", client.HistoryOrder)
}
