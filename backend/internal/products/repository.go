package products

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

func (r *Repository) Create(ctx context.Context, p *Product) error {
	query := `
	   INSERT INTO products (id, name, sku, category, margin_percent, active, image_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.DB.Exec(ctx, query,
		p.ID,
		p.Name,
		p.SKU,
		p.Category,
		p.MarginPercent,
		p.Active,
	)

	return err
}

func (r *Repository) List(ctx context.Context) ([]Product, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id, name, sku, category, margin_percent, active, created_at, image_path
		FROM products
		ORDER BY created_at DESC
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.SKU,
			&p.Category,
			&p.MarginPercent,
			&p.Active,
			&p.CreatedAt,
			&p.ImagePath,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil

}
