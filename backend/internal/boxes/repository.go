package boxes

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(ctx context.Context, b *Box) error {
	query := `
		INSERT INTO boxes (
			id, supplier_name, purchase_date,
			exchange_rate_usd_aoa, shipping_usd,
			customs_tax_usd, other_costs_usd, status
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`

	_, err := r.DB.Exec(ctx, query,
		b.ID,
		b.SupplierName,
		b.PurchaseDate,
		b.ExchangeRateUSDAOA,
		b.ShippingUSD,
		b.CustomsTaxUSD,
		b.OtherCostsUSD,
		b.Status,
	)

	return err
}

func (r *Repository) List(ctx context.Context) ([]Box, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id, supplier_name, purchase_date,
		       exchange_rate_usd_aoa, shipping_usd,
		       customs_tax_usd, other_costs_usd,
		       status, created_at
		FROM boxes
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boxes []Box

	for rows.Next() {
		var b Box
		err := rows.Scan(
			&b.ID,
			&b.SupplierName,
			&b.PurchaseDate,
			&b.ExchangeRateUSDAOA,
			&b.ShippingUSD,
			&b.CustomsTaxUSD,
			&b.OtherCostsUSD,
			&b.Status,
			&b.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		boxes = append(boxes, b)
	}

	return boxes, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Box, error) {
	query := `
		SELECT id, supplier_name, purchase_date,
		       exchange_rate_usd_aoa, shipping_usd,
		       customs_tax_usd, other_costs_usd,
		       status, created_at
		FROM boxes
		WHERE id = $1
	`
	var b Box
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&b.ID,
		&b.SupplierName,
		&b.PurchaseDate,
		&b.ExchangeRateUSDAOA,
		&b.ShippingUSD,
		&b.CustomsTaxUSD,
		&b.OtherCostsUSD,
		&b.Status,
		&b.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (r *Repository) Update(ctx context.Context, b *Box) error {
	query := `
		UPDATE boxes
		SET supplier_name = $1,
		    purchase_date = $2,
		    exchange_rate_usd_aoa = $3,
		    shipping_usd = $4,
		    customs_tax_usd = $5,
		    other_costs_usd = $6,
		    status = $7
		WHERE id = $8
	`
	_, err := r.DB.Exec(ctx, query,
		b.SupplierName,
		b.PurchaseDate,
		b.ExchangeRateUSDAOA,
		b.ShippingUSD,
		b.CustomsTaxUSD,
		b.OtherCostsUSD,
		b.Status,
		b.ID,
	)

	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.Exec(ctx, `
		DELETE FROM boxes
		WHERE id = $1
	`, id)

	return err
}
