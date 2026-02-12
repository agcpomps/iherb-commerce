package admin

import "github.com/labstack/echo/v5"

func RegisterRoutes(e *echo.Group, h *OrdersHandler, p *ProductsHandler, b *BoxesHandler, apikey string) {
	admin := e.Group("/admin", AdminAuthMiddleware())

	admin.POST("/products", p.CreateProduct)
	admin.POST("/boxes", b.CreateBoxes)

	admin.GET("/orders", h.ListOrders)
	admin.POST("/orders/:id/confirm-payment", h.ConfirmPayment)
}
