package sales

import (
	"context"
	"errors"

	"github.com/agcpomps/iherb-commerce/internal/batches"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	DB          *pgxpool.Pool
	BatchesRepo *batches.Repository
}

type SaleItem struct {
	BatchID  uuid.UUID
	Quantity int
	CostAOA  float64
}

func (s *Service) SellProductFIFO(
	ctx context.Context,
	tx pgx.Tx,
	productID string,
	quatity int,
) ([]SaleItem, error) {
	sold, err := s.PreviewProductFIFO(ctx, tx, productID, quatity)
	if err != nil {
		return nil, err
	}

	for _, item := range sold {
		if err := s.BatchesRepo.DecreaseStock(ctx, tx, item.BatchID, item.Quantity); err != nil {
			return nil, err
		}
	}

	return sold, nil
}

func (s *Service) PreviewProductFIFO(
	ctx context.Context,
	tx pgx.Tx,
	productID string,
	quatity int,
) ([]SaleItem, error) {
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

	return sold, nil
}

func (s *Service) ExpireBatchesWithTx(
	ctx context.Context,
	tx pgx.Tx,
) (int, error) {
	return s.BatchesRepo.ExpireOldBatches(ctx, tx)
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

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
