package batches

import (
	"context"

	"github.com/jackc/pgx/v5"
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

func (r *Repository) GetFIFOForProduct(
	ctx context.Context,
	tx pgx.Tx,
	productID string,
) ([]Bacth, error) {
	rows, err := tx.Query(ctx, `
		SELECT id, quantity_available, landed_cost_aoa
		FROM batches
		WHERE product_id = $1
		  AND status = 'active'
		  AND quantity_available > 0
		  AND expiration_date >= CURRENT_DATE
		ORDER BY expiration_date ASC, created_at ASC
		FOR UPDATE`, productID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batches []Bacth
	for rows.Next() {
		var b Bacth
		if err := rows.Scan(&b.ID, &b.QuantityAvailable, &b.LandedCostAOA); err != nil {
			return nil, err
		}
		batches = append(batches, b)
	}

	return batches, nil
}

func (r *Repository) DecreaseStock(
	ctx context.Context,
	tx pgx.Tx,
	batchID string,
	quantity int,
) error {
	_, err := tx.Exec(ctx, `
		UPDATE batches
		SET quantity_available = quantity_available - $1
		WHERE id = $2
		  AND quantity_available >= $1
	`, quantity, batchID)

	return err
}

func (r *Repository) ExpireOldBatches(
	ctx context.Context,
	tx pgx.Tx,
) (int, error) {
	cmd, err := tx.Exec(ctx, `
	UPDATE batches
	SET status = 'expired',
		quantity_available = 0
	WHERE expiration_date < CURRENT_DATE
	  AND status != 'expired'
`)
	if err != nil {
		return 0, err
	}

	return int(cmd.RowsAffected()), err
}
