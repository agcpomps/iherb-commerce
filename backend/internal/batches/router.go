package batches

import "github.com/labstack/echo/v5"

func RegisterRoute(e *echo.Group, h *Handler) {
	e.POST("/batches", h.Create)
}
