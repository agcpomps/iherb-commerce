package products

import "time"

type Product struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	SKU           string    `json:"sku"`
	Category      string    `json:"category"`
	Description   *string   `json:"description"`
	MarginPercent float64   `json:"margin_percent"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"created_at"`
	ImagePath     string    `json:"image_path"`
}

type StoreProduct struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Description       *string `json:"description"`
	ImagePath         string  `json:"image_path"`
	PriceAOA          float64 `json:"price_aoa"`
	AvailableQuantity int     `json:"available_quantity"`
	InStock           bool    `json:"in_stock"`
}
