package orders

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

// Create the order
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, order *Order) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO orders (
			id,
			total_amount_aoa,
			payment_method,
			status
		) VALUES ($1, $2, $3, $4)
	`,
		order.ID,
		order.TotalAmountAOA,
		order.PaymentMethod,
		order.Status,
	)
	return err
}

// criar order item
func (r *Repository) CreateItem(ctx context.Context, tx pgx.Tx, item *OrderItem) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO order_items (
			id,
			order_id,
			product_id,
			product_name,
			batch_id,
			quantity,
			unit_price_aoa,
			total_price_aoa
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`,
		item.ID,
		item.OrderID,
		item.ProductID,
		item.ProductName,
		item.BatchID,
		item.Quantity,
		item.UnitPriceAOA,
		item.TotalPriceAOA,
	)

	return err
}

// Get order by id (useful for admin/confirmation)
func (r *Repository) GetByID(ctx context.Context, db pgx.Tx, orderID string) (*Order, error) {
	var o Order

	err := db.QueryRow(ctx, `
	SELECT
		id,
		total_amount_aoa,
		payment_method,
		status,
		created_at
	FROM orders
	WHERE id = $1
`, orderID).Scan(
		&o.ID,
		&o.TotalAmountAOA,
		&o.PaymentMethod,
		&o.Status,
		&o.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &o, nil
}

// list items from an order
func (r *Repository) ListItems(
	ctx context.Context,
	db pgx.Tx,
	orderID string,
) ([]OrderItem, error) {

	rows, err := db.Query(ctx, `
		SELECT
			id,
			order_id,
			product_id,
			product_name,
			batch_id,
			quantity,
			unit_price_aoa,
			total_price_aoa
		FROM order_items
		WHERE order_id = $1
	`, orderID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem

	for rows.Next() {
		var i OrderItem
		if err := rows.Scan(
			&i.ID,
			&i.OrderID,
			&i.ProductID,
			&i.ProductName,
			&i.BatchID,
			&i.Quantity,
			&i.UnitPriceAOA,
			&i.TotalPriceAOA,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, nil
}

func (r *Repository) MarkAsPaid(
	ctx context.Context,
	tx pgx.Tx,
	orderID string,
) error {
	cmd, err := tx.Exec(ctx, `
	UPDATE orders
	SET status = 'paid'
	WHERE id = $1 AND status = 'pending_payment'
`, orderID)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("order not found or pending")
	}

	return nil
}

func (r *Repository) List(
	ctx context.Context,
	db pgx.Tx,
	status string,
) ([]Order, error) {
	query := `
	SELECT id, total_amount_aoa, payment_method, status, created_at
	FROM orders
`
	args := []interface{}{}
	if status != "" {
		query += " WHERE status = $1"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []Order

	for rows.Next() {
		var o Order
		if err := rows.Scan(
			&o.ID,
			&o.TotalAmountAOA,
			&o.PaymentMethod,
			&o.Status,
			&o.CreatedAt,
		); err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}

	return orders, nil
}
