package admin

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/agcpomps/iherb-commerce/internal/products"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

type ProductsHandler struct {
	DB          *pgxpool.Pool
	ProductRepo *products.Repository
}

const maxImageSize = 5 * 1024 * 1024

func (h *ProductsHandler) ListProducts(c *echo.Context) error {
	ctx := c.Request().Context()

	items, err := h.ProductRepo.List(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to list products",
		})
	}

	return c.JSON(http.StatusOK, items)
}

func (h *ProductsHandler) CreateProduct(c *echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		Name          string  `json:"name"`
		SKU           string  `json:"sku"`
		Category      string  `json:"category"`
		Description   *string `json:"description"`
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
		Description:   req.Description,
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

func (h *ProductsHandler) UpdateProduct(c *echo.Context) error {
	ctx := c.Request().Context()
	productID := c.Param("id")
	if productID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing product id",
		})
	}

	current, err := h.ProductRepo.GetByID(ctx, productID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "product not found",
		})
	}

	var req struct {
		Name          string  `json:"name"`
		SKU           string  `json:"sku"`
		Category      string  `json:"category"`
		Description   *string `json:"description"`
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

	current.Name = req.Name
	current.SKU = req.SKU
	current.Category = req.Category
	current.Description = req.Description
	current.MarginPercent = req.MarginPercent
	current.ImagePath = req.ImagePath
	current.Active = req.Active

	if err := h.ProductRepo.Update(ctx, current); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update product",
		})
	}

	return c.JSON(http.StatusOK, current)
}

func (h *ProductsHandler) DeleteProduct(c *echo.Context) error {
	ctx := c.Request().Context()
	productID := c.Param("id")
	if productID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing product id",
		})
	}

	if _, err := h.ProductRepo.GetByID(ctx, productID); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "product not found",
		})
	}

	if err := h.ProductRepo.Deactivate(ctx, productID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to delete product",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ProductsHandler) UploadProductImage(c *echo.Context) error {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "image file is required",
		})
	}

	if fileHeader.Size > maxImageSize {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "image too large (max 5MB)",
		})
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid image format (use jpg, jpeg, png, or webp)",
		})
	}

	src, err := fileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to open image",
		})
	}
	defer src.Close()

	uploadDir := filepath.Join("uploads", "products")
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to prepare upload directory",
		})
	}

	filename := uuid.NewString() + ext
	dstPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to save image",
		})
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to write image",
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"image_path": "/uploads/products/" + filename,
	})
}
