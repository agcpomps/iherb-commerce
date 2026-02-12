package batches

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

const CUSTOMS_RATE = 0.16

type Handler struct {
	Repo      *Repository
	BoxGetter func(ctx *echo.Context, boxID string) (float64, error)
}

func NewHandler(repo *Repository, boxGetter func(ctx *echo.Context, boxID string) (float64, error)) *Handler {
	return &Handler{Repo: repo, BoxGetter: boxGetter}
}

func (h *Handler) Create(c *echo.Context) error {
	var input struct {
		ProductID        string    `json:"product_id"`
		BoxID            string    `json:"box_id"`
		QuantityReceived int       `json:"quantity_received"`
		ExpirationDate   time.Time `json:"expiration_date"`
		UnitPriceUSD     float64   `json:"unit_price_usd"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	productID, err := uuid.Parse(input.ProductID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid product_id",
		})
	}

	boxID, err := uuid.Parse(input.BoxID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid box_id",
		})
	}

	exchangeRate, err := h.BoxGetter(c, input.BoxID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	customs := input.UnitPriceUSD * CUSTOMS_RATE
	landedUSD := input.UnitPriceUSD + customs
	landedAOA := landedUSD * exchangeRate

	batch := Bacth{
		ID:                uuid.New(),
		ProductID:         productID,
		BoxID:             boxID,
		QuantityReceived:  input.QuantityReceived,
		QuantityAvailable: input.QuantityReceived,
		ExpirationDate:    input.ExpirationDate,
		LandedCostUSD:     landedUSD,
		LandedCostAOA:     landedAOA,
		Status:            "active",
	}

	if err := h.Repo.Create(c.Request().Context(), &batch); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, batch)
}
