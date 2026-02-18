package admin

import (
	"net/http"
	"time"

	"github.com/agcpomps/iherb-commerce/internal/boxes"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type BoxesHandler struct {
	BoxesRepo *boxes.Repository
}

func (h *BoxesHandler) ListBoxes(c *echo.Context) error {
	ctx := c.Request().Context()

	items, err := h.BoxesRepo.List(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to list boxes",
		})
	}

	return c.JSON(http.StatusOK, items)
}

func (h *BoxesHandler) CreateBoxes(c *echo.Context) error {
	var req struct {
		SupplierName       string  `json:"supplier_name"`
		PurchaseDate       string  `json:"purchase_date"`
		ExchangeRateUSDAOA float64 `json:"exchange_rate_usd_aoa"`
		ShippingUSD        float64 `json:"shipping_usd"`
		CustomsTaxUSD      float64 `json:"customs_tax_usd"`
		OtherCostsUSD      float64 `json:"other_costs_usd"`
		Status             string  `json:"status"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	if req.SupplierName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "supplier_name required",
		})
	}

	date, err := time.Parse("2006-01-02", req.PurchaseDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid purchase_date format",
		})
	}

	if req.Status != "in_transit" && req.Status != "arrived" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid status",
		})
	}

	box := &boxes.Box{
		ID:                 uuid.NewString(),
		SupplierName:       req.SupplierName,
		PurchaseDate:       date,
		ExchangeRateUSDAOA: req.ExchangeRateUSDAOA,
		ShippingUSD:        req.ShippingUSD,
		CustomsTaxUSD:      req.CustomsTaxUSD,
		OtherCostsUSD:      req.OtherCostsUSD,
		Status:             req.Status,
	}

	ctx := c.Request().Context()

	if err := h.BoxesRepo.Create(ctx, box); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create box",
		})
	}

	return c.JSON(http.StatusCreated, box)

}

func (h *BoxesHandler) UpdateBox(c *echo.Context) error {
	ctx := c.Request().Context()
	boxID := c.Param("id")
	if boxID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing box id",
		})
	}

	current, err := h.BoxesRepo.GetByID(ctx, boxID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "box not found",
		})
	}

	var req struct {
		SupplierName       string  `json:"supplier_name"`
		PurchaseDate       string  `json:"purchase_date"`
		ExchangeRateUSDAOA float64 `json:"exchange_rate_usd_aoa"`
		ShippingUSD        float64 `json:"shipping_usd"`
		CustomsTaxUSD      float64 `json:"customs_tax_usd"`
		OtherCostsUSD      float64 `json:"other_costs_usd"`
		Status             string  `json:"status"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	if req.SupplierName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "supplier_name required",
		})
	}

	date, err := time.Parse("2006-01-02", req.PurchaseDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid purchase_date format",
		})
	}

	if req.Status != "in_transit" && req.Status != "arrived" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid status",
		})
	}

	current.SupplierName = req.SupplierName
	current.PurchaseDate = date
	current.ExchangeRateUSDAOA = req.ExchangeRateUSDAOA
	current.ShippingUSD = req.ShippingUSD
	current.CustomsTaxUSD = req.CustomsTaxUSD
	current.OtherCostsUSD = req.OtherCostsUSD
	current.Status = req.Status

	if err := h.BoxesRepo.Update(ctx, current); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update box",
		})
	}

	return c.JSON(http.StatusOK, current)
}

func (h *BoxesHandler) DeleteBox(c *echo.Context) error {
	ctx := c.Request().Context()
	boxID := c.Param("id")
	if boxID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing box id",
		})
	}

	if _, err := h.BoxesRepo.GetByID(ctx, boxID); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "box not found",
		})
	}

	if err := h.BoxesRepo.Delete(ctx, boxID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to delete box",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
