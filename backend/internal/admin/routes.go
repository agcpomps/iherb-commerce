package admin

import "github.com/labstack/echo/v5"

func RegisterRoutes(e *echo.Group, h *OrdersHandler, p *ProductsHandler, b *BoxesHandler, bh *BatchesHandler) {
	admin := e.Group("/admin", AdminAuthMiddleware())

	admin.GET("/products", p.ListProducts)
	admin.POST("/products", p.CreateProduct)
	admin.POST("/uploads/products", p.UploadProductImage)
	admin.PUT("/products/:id", p.UpdateProduct)
	admin.DELETE("/products/:id", p.DeleteProduct)
	admin.GET("/boxes", b.ListBoxes)
	admin.POST("/boxes", b.CreateBoxes)
	admin.PUT("/boxes/:id", b.UpdateBox)
	admin.DELETE("/boxes/:id", b.DeleteBox)
	admin.POST("/batches", bh.CreateBatch)

	admin.GET("/orders", h.ListOrders)
	admin.GET("/orders/:id", h.GetOrder)
	admin.POST("/orders/:id/confirm-payment", h.ConfirmPayment)
	admin.POST("/orders/:id/cancel", h.CancelOrder)
}
