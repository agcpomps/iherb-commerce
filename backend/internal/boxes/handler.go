package boxes

import (
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
	var input Box
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	input.ID = uuid.NewString()
	if input.SupplierName == "" {
		input.SupplierName = "iherb"
	}

	if input.Status == "" {
		input.Status = "in_transit"
	}

	if err := h.Repo.Create(c.Request().Context(), &input); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, input)
}

func (h *Handler) List(c *echo.Context) error {
	boxes, err := h.Repo.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, boxes)
}

func (h *Handler) GetByID(c *echo.Context) error {
	id := c.Param("id")

	box, err := h.Repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, box)
}

func (h *Handler) Update(c *echo.Context) error {
	id := c.Param("id")

	var input Box

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	input.ID = id

	if err := h.Repo.Update(c.Request().Context(), &input); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, input)
}
