package links

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/SKjustSK/alru-url-shortener/backend/internal/database"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/models"
	"github.com/SKjustSK/alru-url-shortener/backend/pkg/base62"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type CreateLinkRequest struct {
	LongURL   string     `json:"long_url"`
	ShortCode *string    `json:"short_code,omitempty"` // Used as Custom Alias if provided
	ExpiresOn *time.Time `json:"expires_on,omitempty"`
}

type CreateLinkResponse struct {
	LongURL   string    `json:"long_url"`
	ShortURL  string    `json:"short_url"`
	ExpiresOn time.Time `json:"expires_on"`
}

func CreateLink(c *echo.Context) error {
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
		if expiryTime.Before(time.Now().UTC()) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Expiration date must be in the future",
			})
		}
	} else {
		// Default to 24 hours if not provided
		expiryTime = time.Now().UTC().Add(24 * time.Hour)
	}

	// 3. Setup Initial Link Record
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*models.JWTCustomClaims)

	newLink := models.Link{
		UserID:    claims.UserID,
		LongURL:   req.LongURL,
		ExpiresAt: expiryTime,
	}

	// 4. Link insertiong
	isCustom := req.ShortCode != nil && *req.ShortCode != ""
	newLink.IsCustom = isCustom

	if isCustom {
		// --- CUSTOM LINK ---
		newLink.ShortCode = *req.ShortCode

		// Attempt to insert directly. PostgreSQL's composite unique index will block duplicates.
		if err := database.DB.Create(&newLink).Error; err != nil {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "This custom alias is already in use. Please choose another.",
			})
		}
	} else {
		// --- GENERATED LINK ---
		newLink.ShortCode = "PENDING" // Safe placeholder

		// Insert first to let PostgreSQL generate the auto-incrementing LinkID
		if err := database.DB.Create(&newLink).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to initialize link"})
		}

		// Because it uses the unique LinkID, collisions are mathematically impossible.
		newLink.ShortCode = base62.Encode(uint64(newLink.LinkID))

		// Update the row with the newly generated secure code
		if err := database.DB.Model(&newLink).Update("short_code", newLink.ShortCode).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to finalize short code"})
		}
	}

	// 5. Redis Caching
	ctx := c.Request().Context()
	redisTTL := min(24*time.Hour, time.Until(expiryTime))

	// Namespace the Redis keys to prevent overlapping
	redisKey := newLink.ShortCode
	if newLink.IsCustom {
		redisKey = "custom:" + newLink.ShortCode
	} else {
		redisKey = "gen:" + newLink.ShortCode
	}

	_ = database.RedisDB.Set(ctx, redisKey, newLink.LongURL, redisTTL).Err()

	// 6. Final Response Construction
	baseURL := os.Getenv("BACKEND_URL")
	if baseURL == "" {
		baseURL = "http://localhost:1323"
	}

	// Apply the route prefixing if it is a custom link
	prefix := "/"
	if newLink.IsCustom {
		prefix = "/c/"
	}

	return c.JSON(http.StatusCreated, CreateLinkResponse{
		ShortURL:  baseURL + prefix + newLink.ShortCode,
		LongURL:   newLink.LongURL,
		ExpiresOn: expiryTime,
	})
}
