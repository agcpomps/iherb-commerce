package batches

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

func (r *Repository) Create(ctx context.Context, b *Bacth) error {
	query := `
		INSERT INTO batches (
			id, product_id, box_id,
			quantity_received, quantity_available,
			expiration_date, landed_cost_usd,
			landed_cost_aoa, status
		)
		VALUES ($1,$2,$3,$4,$4,$5,$6,$7,$8)
	`
	_, err := r.DB.Exec(ctx, query,
		b.ID,
		b.ProductID,
		b.BoxID,
		b.QuantityReceived,
		b.ExpirationDate,
		b.LandedCostUSD,
		b.LandedCostAOA,
		b.Status,
	)
	return err
}
