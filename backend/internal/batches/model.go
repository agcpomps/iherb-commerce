package batches

import (
	"time"

	"github.com/google/uuid"
)

type Bacth struct {
	ID                uuid.UUID `json:"id"`
	ProductID         uuid.UUID `json:"product_id"`
	BoxID             uuid.UUID `json:"box_id"`
	QuantityReceived  int       `json:"quantity_received"`
	QuantityAvailable int       `json:"quantity_available"`
	ExpirationDate    time.Time `json:"expiration_date"`
	LandedCostUSD     float64   `json:"landed_cost_usd"`
	LandedCostAOA     float64   `json:"landed_cost_aoa"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}
