package products

import "time"

type Product struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	SKU           string    `json:"sku"`
	Category      string    `json:"category"`
	MarginPercent float64   `json:"margin_percent"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"created_at"`
	ImagePath     string    `json:"image_path"`
}
