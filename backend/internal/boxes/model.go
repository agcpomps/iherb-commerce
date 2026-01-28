package boxes

import "time"

type Box struct {
	ID                 string    `json:"id"`
	SupplierName       string    `json:"supplier_name"`
	PurchaseDate       time.Time `json:"purchase_date"`
	ExchangeRateUSDAOA float64   `json:"exchange_rate_usd_aoa"`
	ShippingUSD        float64   `json:"shipping_usd"`
	CustomsTaxUSD      float64   `json:"customs_tax_usd"`
	OtherCostsUSD      float64   `json:"other_costs_usd"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
}
