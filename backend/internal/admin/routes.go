package admin

import "github.com/labstack/echo/v5"

func RegisterRoutes(e *echo.Group, h *OrdersHandler, p *ProductsHandler, apikey string) {
	admin := e.Group("/admin", AdminAuthMiddleware())

	admin.POST("/products", p.CreateProduct)

	admin.GET("/orders", h.ListOrders)
	admin.POST("/orders/:id/confirm-payment", h.ConfirmPayment)
}
