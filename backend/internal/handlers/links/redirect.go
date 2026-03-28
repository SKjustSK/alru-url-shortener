package links

import (
	"net/http"
	"time"

	"github.com/SKjustSK/alru-url-shortener/backend/internal/database"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/models"
	"github.com/labstack/echo/v5"
)

func RedirectLink(c *echo.Context) error {
	shortCode := c.Param("short_code")
	ctx := c.Request().Context()

	// 1. Check Redis (The Fast Path)
	longURL, err := database.RedisDB.Get(ctx, shortCode).Result()
	if err == nil {
		// Cache Hit: Redirect immediately
		return c.Redirect(http.StatusFound, longURL)
	}

	// 2. Cache Miss: Query PostgreSQL
	var link models.Link

	err = database.DB.Select("link_id", "long_url", "expires_at").
		Where("short_code = ?", shortCode).
		Order("expires_at DESC").
		First(&link).Error

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Link not found"})
	}

	// 3. Expiration Check
	if time.Now().After(link.ExpiresAt) {
		return c.JSON(http.StatusGone, map[string]string{"error": "This link has expired"})
	}

	// 4. Re-populate Redis
	redisTTL := min(24*time.Hour, time.Until(link.ExpiresAt))

	_ = database.RedisDB.Set(ctx, shortCode, link.LongURL, redisTTL).Err()

	// 5. Final Redirect
	// We use http.StatusFound (302) so browsers don't permanently cache the redirect,
	// ensuring every click hits our server so we can track analytics.
	return c.Redirect(http.StatusFound, link.LongURL)
}
