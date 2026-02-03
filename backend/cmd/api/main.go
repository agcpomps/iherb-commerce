package main

import (
	"log"
	"os"

	"github.com/agcpomps/iherb-commerce/internal/batches"
	"github.com/agcpomps/iherb-commerce/internal/boxes"
	"github.com/agcpomps/iherb-commerce/internal/db"
	"github.com/agcpomps/iherb-commerce/internal/products"
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

	e.GET("api/v1/health", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
		})
	})

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

	api := e.Group("/api/v1")
	products.RegisterRoutes(api, handler)
	boxes.RegisterRoutes(api, boxHandler)
	batches.RegisterRoute(api, batchHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("API running on port", port)

	if err := e.Start("0.0.0.0:" + port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

}
