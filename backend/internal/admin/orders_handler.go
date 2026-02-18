package admin

import (
	"net/http"

	"github.com/agcpomps/iherb-commerce/internal/batches"
	"github.com/agcpomps/iherb-commerce/internal/orders"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

type OrdersHandler struct {
	DB          *pgxpool.Pool
	OrdersRepo  *orders.Repository
	BatchesRepo *batches.Repository
}

func (h *OrdersHandler) ConfirmPayment(c *echo.Context) error {
	ctx := c.Request().Context()
	orderID := c.Param("id")

	tx, err := h.DB.Begin(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to start transaction",
		})
	}

	defer tx.Rollback(ctx)

	order, err := h.OrdersRepo.GetByID(ctx, tx, orderID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "order not found",
		})
	}
	if order.Status != "pending_payment" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "order is not pending payment",
		})
	}

	items, err := h.OrdersRepo.ListItems(ctx, tx, orderID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "failed to load order items",
		})
	}
	if len(items) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "order has no items",
		})
	}

	for _, item := range items {
		if err := h.BatchesRepo.DecreaseStock(ctx, tx, item.BatchID, item.Quantity); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "insufficient stock to confirm this order",
			})
		}
	}

	err = h.OrdersRepo.MarkAsPaid(ctx, tx, orderID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to commit",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"order_id": orderID,
		"status":   "paid",
	})
}

func (h *OrdersHandler) GetOrder(c *echo.Context) error {
	ctx := c.Request().Context()
	orderID := c.Param("id")

	tx, err := h.DB.Begin(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to start transaction",
		})
	}
	defer tx.Rollback(ctx)

	order, err := h.OrdersRepo.GetByID(ctx, tx, orderID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "order not found",
		})
	}

	items, err := h.OrdersRepo.ListItems(ctx, tx, orderID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to load order items",
		})
	}
	order.Items = items

	if err := tx.Commit(ctx); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to commit",
		})
	}

	return c.JSON(http.StatusOK, order)
}

func (h *OrdersHandler) CancelOrder(c *echo.Context) error {
	ctx := c.Request().Context()
	orderID := c.Param("id")

	tx, err := h.DB.Begin(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to start transaction",
		})
	}
	defer tx.Rollback(ctx)

	if err := h.OrdersRepo.MarkAsCancelled(ctx, tx, orderID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to commit",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"order_id": orderID,
		"status":   "cancelled",
	})
}

func (h *OrdersHandler) ListOrders(c *echo.Context) error {
	ctx := c.Request().Context()

	status := c.QueryParam("status")

	tx, err := h.DB.Begin(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to start transaction",
		})
	}

	defer tx.Rollback(ctx)

	orders, err := h.OrdersRepo.List(ctx, tx, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to list orders",
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to commit",
		})
	}

	return c.JSON(http.StatusOK, orders)
}
