package admin

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
)

func AdminAuthMiddleware() echo.MiddlewareFunc {
	expectedHash := os.Getenv("ADMIN_API_KEY_HASH")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			apiKey := c.Request().Header.Get("X-API-KEY")
			if apiKey == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing api key",
				})
			}

			hash := sha256.Sum256([]byte(apiKey))
			hashHex := hex.EncodeToString(hash[:])

			if hashHex != expectedHash {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "unauthorized",
				})
			}

			return next(c)
		}
	}
}
