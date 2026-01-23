package products

import "github.com/labstack/echo/v5"

func RegisterRoutes(e *echo.Group, h *Handler) {
	e.POST("/products", h.Create)
	e.GET("/products", h.List)
}
