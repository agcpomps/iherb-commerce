package admin

import (
	"net/http"

	"github.com/agcpomps/iherb-commerce/internal/products"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

type ProductsHandler struct {
	DB          *pgxpool.Pool
	ProductRepo *products.Repository
}

func (h *ProductsHandler) CreateProduct(c *echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		Name          string  `json:"name"`
		SKU           string  `json:"sku"`
		Category      string  `json:"category"`
		MarginPercent float64 `json:"margin_percent"`
		ImagePath     string  `json:"image_path"`
		Active        bool    `json:"active"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	if req.Name == "" || req.SKU == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "name and sku required",
		})
	}

	product := &products.Product{
		ID:            uuid.NewString(),
		Name:          req.Name,
		SKU:           req.SKU,
		Category:      req.Category,
		MarginPercent: req.MarginPercent,
		Active:        req.Active,
		ImagePath:     req.ImagePath,
	}

	if err := h.ProductRepo.Create(ctx, product); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create product",
		})
	}

	return c.JSON(http.StatusCreated, product)
}
