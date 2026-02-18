package main

import (
	"log"
	"net/http"
	"os"

	"github.com/agcpomps/iherb-commerce/internal/admin"
	"github.com/agcpomps/iherb-commerce/internal/batches"
	"github.com/agcpomps/iherb-commerce/internal/boxes"
	"github.com/agcpomps/iherb-commerce/internal/db"
	"github.com/agcpomps/iherb-commerce/internal/orders"
	"github.com/agcpomps/iherb-commerce/internal/products"
	"github.com/agcpomps/iherb-commerce/internal/sales"
	"github.com/agcpomps/iherb-commerce/internal/store"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {

	pool, err := db.NewPool()
	if err != nil {
		log.Fatal(err)
	}
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080", "http://localhost:4321"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderXRequestedWith,
			"X-API-KEY",
		},
	}))

	e.GET("/api/v1/health", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
		})
	})
	e.Static("/uploads", "uploads")

	repo := products.NewRepository(pool)
	handler := products.NewHandler(repo)
	boxRepo := boxes.NewRepository(pool)
	boxHandler := boxes.NewHandler(boxRepo)
	batchRepo := batches.NewRepository(pool)
	batchHandler := batches.NewHandler(batchRepo, func(c *echo.Context, boxID string) (float64, error) {
		var rate float64
		err := pool.QueryRow(
			c.Request().Context(),
			"SELECT exchange_rate_usd_aoa FROM boxes WHERE id = $1",
			boxID,
		).Scan(&rate)

		return rate, err
	})
	orderRepo := orders.NewRepository()
	salesService := &sales.Service{
		DB:          pool,
		BatchesRepo: batchRepo,
	}
	storeHandler := store.NewHandler(pool, repo, orderRepo, salesService)
	api := e.Group("/api/v1")
	products.RegisterRoutes(api, handler)
	boxes.RegisterRoutes(api, boxHandler)
	batches.RegisterRoute(api, batchHandler)

	// store register
	store.RegisterRoutes(api, storeHandler)
	store.RegisterRoutes(e.Group(""), storeHandler)

	// admin
	adminOrdersHandler := &admin.OrdersHandler{
		DB:          pool,
		OrdersRepo:  orderRepo,
		BatchesRepo: batchRepo,
	}

	productsHandler := &admin.ProductsHandler{
		DB:          pool,
		ProductRepo: repo,
	}

	boxesHandler := &admin.BoxesHandler{
		BoxesRepo: boxRepo,
	}

	batchesHandler := &admin.BatchesHandler{
		BatchHandler: batchHandler,
	}

	admin.RegisterRoutes(api, adminOrdersHandler, productsHandler, boxesHandler, batchesHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("API running on port", port)

	if err := e.Start("0.0.0.0:" + port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

}
