package batches

import "time"

type Bacth struct {
	ID                string    `json:"id"`
	ProductID         string    `json:"product_id"`
	BoxID             string    `json:"box_id"`
	QuantityReceived  int       `json:"quantity_received"`
	QuantityAvailable int       `json:"quantity_available"`
	ExpirationDate    time.Time `json:"expiration_date"`
	LandedCostUSD     float64   `json:"landed_cost_usd"`
	LandedCostAOA     float64   `json:"landed_cost_aoa"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}
