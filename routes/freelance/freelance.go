package freelance

import (
	"api-smart-room/api/freelance"

	"github.com/labstack/echo/v4"
)

func FreelanceSubRoute(group *echo.Group) {
	group.GET("/offerings", freelance.OfferingList)
	group.GET("/offerings/:id_order", freelance.OfferingDetail)
	group.PATCH("/offerings/:id_order/confirm", freelance.ConfirmOffering)
	group.PATCH("/offerings/:id_order/reject", freelance.RejectOffering)

	group.GET("/offerings/:id_order/coordinate-both", freelance.GetCoordinateBoth)
	group.POST("/offerings/:id_order/arrangement/task", freelance.AddTask)
	group.DELETE("/offerings/:id_order/arrangement/task/:id_task", freelance.DeleteTask)
	group.POST("/offerings/:id_order/arrangement", freelance.ArrangeOffering)
	group.GET("/offerings/:id_order/arrangement", freelance.GetArrangement)

	group.GET("/offerings/:id_order/status", freelance.RefreshStatus)

	group.GET("/history", freelance.HistoriOffering)

	group.GET("/me", freelance.GetProfile)
	group.PATCH("/me/update-address", freelance.UpdateAddress)
}
