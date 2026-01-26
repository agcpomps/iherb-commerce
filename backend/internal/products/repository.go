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
		p.ImagePath,
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

func (r *Repository) GetByID(ctx context.Context, id string) (*Product, error) {
	query := `
		SELECT id, name, sku, category, margin_percent, active, created_at, image_path
		FROM products
		WHERE id = $1
	`
	var p Product

	err := r.DB.QueryRow(ctx, query, id).Scan(
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

	return &p, nil

}

func (r *Repository) Update(ctx context.Context, p *Product) error {
	query := `
		UPDATE products
		SET name = $1,
		    sku = $2,
		    category = $3,
		    margin_percent = $4,
		    active = $5
			image_path = &6
		WHERE id = $7
	`
	_, err := r.DB.Exec(ctx, query,
		p.Name,
		p.SKU,
		p.Category,
		p.MarginPercent,
		p.Active,
		p.ImagePath,
		p.ID,
	)

	return err
}

func (r *Repository) Deactivate(ctx context.Context, id string) error {
	query := `
		UPDATE products
		SET active = false
		WHERE id = $1
	`
	_, err := r.DB.Exec(ctx, query, id)
	return err
}
