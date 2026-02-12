package orders

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID             string      `json:"id"`
	TotalAmountAOA float64     `json:"total_amount_aoa"`
	PaymentMethod  string      `json:"payment_method"`
	Status         string      `json:"status"`
	CreatedAt      time.Time   `json:"created_at"`
	Items          []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID            uuid.UUID `json:"id"`
	OrderID       string    `json:"order_id"`
	ProductID     uuid.UUID `json:"product_id"`
	ProductName   string    `json:"product_name"`
	BatchID       uuid.UUID `json:"batch_id"`
	Quantity      int       `json:"quantity"`
	UnitPriceAOA  float64   `json:"unit_price_aoa"`
	TotalPriceAOA float64   `json:"total_price_aoa"`
}
