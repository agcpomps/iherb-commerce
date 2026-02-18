package admin

import (
	"net/http"

	"github.com/agcpomps/iherb-commerce/internal/batches"
	"github.com/labstack/echo/v5"
)

type BatchesHandler struct {
	BatchHandler *batches.Handler
}

func (h *BatchesHandler) CreateBatch(c *echo.Context) error {
	if h.BatchHandler == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "batch handler not configured",
		})
	}

	return h.BatchHandler.Create(c)
}
