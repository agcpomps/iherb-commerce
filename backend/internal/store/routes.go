package store

import "github.com/labstack/echo/v5"

func RegisterRoutes(e *echo.Group, h *Handler) {
	e.POST("/store/checkout", h.Checkout)
	e.GET("/store/products", h.ListProducts)
}
