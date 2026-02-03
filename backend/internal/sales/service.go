package sales

import (
	"context"
	"errors"

	"github.com/agcpomps/iherb-commerce/internal/batches"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	DB          *pgxpool.Pool
	BatchesRepo *batches.Repository
}

type SaleItem struct {
	BatchID  string
	Quantity int
	CostAOA  float64
}

func (s *Service) SellProductFIFO(
	ctx context.Context,
	productID string,
	quatity int,
) ([]SaleItem, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	batches, err := s.BatchesRepo.GetFIFOForProduct(ctx, tx, productID)
	if err != nil {
		return nil, err
	}

	remaining := quatity
	var sold []SaleItem

	for _, b := range batches {
		if remaining == 0 {
			break
		}

		toSell := min(b.QuantityAvailable, remaining)

		if err := s.BatchesRepo.DecreaseStock(ctx, tx, b.ID, toSell); err != nil {
			return nil, err
		}

		sold = append(sold, SaleItem{
			BatchID:  b.ID,
			Quantity: toSell,
			CostAOA:  b.LandedCostAOA,
		})

		remaining -= toSell
	}

	if remaining > 0 {
		return nil, errors.New("stock insuficiente")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return sold, nil
}

func (s *Service) ExpireBatches(ctx context.Context) (int, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return 0, err
	}

	defer tx.Rollback(ctx)

	affected, err := s.BatchesRepo.ExpireOldBatches(ctx, tx)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return affected, nil
}
