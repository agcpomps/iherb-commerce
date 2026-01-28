package boxes

import "github.com/labstack/echo/v5"

func RegisterRoutes(e *echo.Group, h *Handler) {
	e.POST("/boxes", h.Create)
	e.GET("/boxes", h.List)
	e.GET("/boxes/:id", h.GetByID)
	e.PUT("/boxes/:id", h.Update)
}
