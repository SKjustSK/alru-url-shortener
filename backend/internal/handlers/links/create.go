package links

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/SKjustSK/alru-url-shortener/backend/internal/database"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/models"
	"github.com/SKjustSK/alru-url-shortener/backend/pkg/base62"
	"github.com/labstack/echo/v5"
)

type CreateLinkRequest struct {
	LongURL   string     `json:"long_url"`
	ExpiresOn *time.Time `json:"expires_on,omitempty"`
}

type CreateLinkResponse struct {
	ShortURL  string    `json:"short_url"`
	LongURL   string    `json:"long_url"`
	ExpiresOn time.Time `json:"expires_on"`
}

func CreateLink(c *echo.Context) error { // Removed the * from echo.Context
	req := new(CreateLinkRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON request",
		})
	}

	// 1. Validate LongURL
	if req.LongURL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL is required"})
	}
	parsedURL, err := url.ParseRequestURI(req.LongURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "A valid URL including http/https is required",
		})
	}

	// 2. Determine Expiration
	var expiryTime time.Time
	if req.ExpiresOn != nil {
		expiryTime = *req.ExpiresOn
		if expiryTime.Before(time.Now()) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Expiration date must be in the future",
			})
		}
	} else {
		expiryTime = time.Now().Add(24 * time.Hour)
	}

	// 3. Setup Link Record with Placeholder
	// Since ShortCode is NOT NULL, we put a temporary value.
	newLink := models.Link{
		LongURL:   req.LongURL,
		ShortCode: "pending",
		ExpiresAt: expiryTime,
	}

	// 4. Create Record to generate the auto-increment LinkID
	if err := database.DB.Create(&newLink).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database Save Failed"})
	}

	// 5. Generate ShortCode via Base62 Counter
	newLink.ShortCode = base62.Encode(uint64(newLink.LinkID))

	// 6. Persist the generated code back to the DB
	if err := database.DB.Model(&newLink).Update("ShortCode", newLink.ShortCode).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update short code"})
	}

	// 7. Redis Caching (Capped at 24h)
	redisTTL := min(24*time.Hour, time.Until(expiryTime))
	ctx := c.Request().Context()

	// Error is ignored here so the user still gets their link if Redis is momentarily down
	_ = database.RedisDB.Set(ctx, newLink.ShortCode, newLink.LongURL, redisTTL).Err()

	// 8. Final Response
	baseURL := os.Getenv("BACKEND_URL")
	if baseURL == "" {
		baseURL = "http://localhost:1323"
	}

	return c.JSON(http.StatusCreated, CreateLinkResponse{
		ShortURL:  baseURL + "/" + newLink.ShortCode,
		LongURL:   newLink.LongURL,
		ExpiresOn: expiryTime,
	})
}
