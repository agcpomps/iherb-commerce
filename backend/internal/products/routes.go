package products

import "github.com/labstack/echo/v5"

func RegisterRoutes(e *echo.Group, h *Handler) {
	e.POST("/products", h.Create)
	e.GET("/products", h.List)
	e.GET("/products/:id", h.GetByID)
	e.PUT("/products/:id", h.Update)
	e.DELETE("/products/:id", h.Deactivate)
}
