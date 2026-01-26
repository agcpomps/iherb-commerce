package products

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	Repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{Repo: repo}
}

func (h *Handler) Create(c *echo.Context) error {
	var input Product

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	input.ID = uuid.NewString()
	input.Active = true

	if err := h.Repo.Create(c.Request().Context(), &input); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, input)
}

func (h *Handler) List(c *echo.Context) error {
	products, err := h.Repo.List(c.Request().Context())
	fmt.Printf("the products: %v", products)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, products)
}

func (h *Handler) GetByID(c *echo.Context) error {
	id := c.Param("id")

	product, err := h.Repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "product not found",
		})
	}

	return c.JSON(http.StatusOK, product)
}

func (h *Handler) Update(c *echo.Context) error {
	id := c.Param("id")

	var input Product
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	input.ID = id

	if err := h.Repo.Update(c.Request().Context(), &input); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, input)
}

func (h *Handler) Deactivate(c *echo.Context) error {
	id := c.Param("id")

	if err := h.Repo.Deactivate(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
