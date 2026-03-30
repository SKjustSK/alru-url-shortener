package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/SKjustSK/alru-url-shortener/backend/internal/database"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func LinkRateLimiter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		// 1. Get UserID from JWT context
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*models.JWTCustomClaims)
		userID := claims.UserID

		ctx := context.Background()
		key := fmt.Sprintf("ratelimit:links:%d", userID)

		// 2. Increment the count in Redis
		count, err := database.RedisDB.Incr(ctx, key).Result()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Rate limiter error")
		}

		// 3. Set expiration to 1 hour on the first hit
		if count == 1 {
			database.RedisDB.Expire(ctx, key, 1*time.Hour)
		}

		// 4. Check if limit exceeded
		if count > 5 {
			return echo.NewHTTPError(http.StatusTooManyRequests, "Limit reached: 5 links per hour")
		}

		return next(c)
	}
}
