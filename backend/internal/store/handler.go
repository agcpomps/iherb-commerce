package store

import (
	"fmt"
	"net/http"
	"time"

	"github.com/agcpomps/iherb-commerce/internal/orders"
	"github.com/agcpomps/iherb-commerce/internal/products"
	"github.com/agcpomps/iherb-commerce/internal/sales"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	DB           *pgxpool.Pool
	ProductsRepo *products.Repository
	OrdersRepo   *orders.Repository
	SalesService *sales.Service
}

func NewHandler(db *pgxpool.Pool, productsRepo *products.Repository, ordersRepo *orders.Repository, salesService *sales.Service) *Handler {
	return &Handler{
		DB:           db,
		ProductsRepo: productsRepo,
		OrdersRepo:   ordersRepo,
		SalesService: salesService,
	}
}

type CheckoutItemInput struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CheckoutRequest struct {
	Items        []CheckoutItemInput `json:"items"`
	PaymentMehod string              `json:"payment_method"`
}

func (h *Handler) Checkout(c *echo.Context) error {
	ctx := c.Request().Context()

	var req CheckoutRequest
	if err := c.Bind(&req); err != nil || len(req.Items) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	// validate payment method
	if req.PaymentMehod != "transfer_site" && req.PaymentMehod != "transfer_whatsapp" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid payment method",
		})
	}

	tx, err := h.DB.Begin(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to start transactin",
		})
	}

	defer tx.Rollback(ctx)

	orderID := fmt.Sprintf("ORD-%s", time.Now().Format("20060102-150405"))

	var totalAOA float64
	var orderItems []orders.OrderItem

	for _, item := range req.Items {
		if item.Quantity <= 0 || item.Quantity > 5 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid quantity",
			})
		}

		productUUID, err := uuid.Parse(item.ProductID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid product_id",
			})
		}

		// product + margin
		product, err := h.ProductsRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "product not found",
			})
		}

		// FIFO (baixa stock)
		sold, err := h.SalesService.SellProductFIFO(ctx, tx, item.ProductID, item.Quantity)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		for _, s := range sold {
			unitPrice := s.CostAOA * (1 + (product.MarginPercent / 100))
			itemTotal := unitPrice * float64(s.Quantity)

			orderItems = append(orderItems, orders.OrderItem{
				ID:            uuid.New(),
				OrderID:       orderID,
				ProductID:     productUUID,
				ProductName:   product.Name,
				BatchID:       s.BatchID,
				Quantity:      s.Quantity,
				UnitPriceAOA:  unitPrice,
				TotalPriceAOA: itemTotal,
			})

			totalAOA += itemTotal
		}
	}

	err = h.OrdersRepo.Create(ctx, tx, &orders.Order{
		ID:             orderID,
		TotalAmountAOA: totalAOA,
		PaymentMethod:  req.PaymentMehod,
		Status:         "pending_payment",
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create order",
		})
	}

	for _, item := range orderItems {
		orderItem := item
		err = h.OrdersRepo.CreateItem(ctx, tx, &orderItem)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to create order item",
			})
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to commit",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"order_id":  orderID,
		"status":    "pending_payment",
		"total_aoa": totalAOA,
		"bank_details": map[string]string{
			"bank":         "BAI",
			"account_name": "IHerb Store Benguela",
			"iban":         "AO06XXXXXXX",
		},
		"instructions": []string{
			"Faça a transferência do valor acima",
			"Use o número do pedido como referência",
			"Envie o comprovativo por WhatsApp ou Email",
		},
		"whatsapp_url": buildWhatsAppURL(orderID, totalAOA),
		"email":        "vendas@seudominio.com",
	})

}
