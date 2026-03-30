package links

import (
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/SKjustSK/alru-url-shortener/backend/internal/database"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/models"
	"github.com/SKjustSK/alru-url-shortener/backend/pkg/base62"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type CreateLinkRequest struct {
	LongURL   string     `json:"long_url"`
	ShortCode *string    `json:"short_code,omitempty"`
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

	// 3. Setup Link Record
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*models.JWTCustomClaims)
	newLink := models.Link{
		UserID:    claims.UserID,
		LongURL:   req.LongURL,
		ShortCode: "PENDING", // Safe placeholder
		ExpiresAt: expiryTime,
	}

	// If user provided a custom code (e.g. "promo"), set it now
	if req.ShortCode != nil && *req.ShortCode != "" {
		newLink.ShortCode = *req.ShortCode
	}

	// 4.1 Check Redis for short code
	ctx := c.Request().Context()
	if newLink.ShortCode != "PENDING" {
		exists, err := database.RedisDB.Exists(ctx, newLink.ShortCode).Result()

		if err == nil && exists > 0 {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "Custom short code is currently in use",
			})
		}
	}

	// 4.2 Create Record Safely using a Database Transaction
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// If they want a custom code, check if it's currently active
		if newLink.ShortCode != "PENDING" {
			var count int64
			tx.Model(&models.Link{}).
				Where("short_code = ? AND expires_at > ?", newLink.ShortCode, time.Now().UTC()).
				Count(&count)

			if count > 0 {
				// Trigger a rollback and pass this specific error out
				return errors.New("CODE_ACTIVE")
			}
		}

		// Insert the record. If it succeeds, the transaction commits.
		if err := tx.Create(&newLink).Error; err != nil {
			return err
		}

		return nil
	})

	// Handle the results of the transaction
	if err != nil {
		if err.Error() == "CODE_ACTIVE" {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "This short code is currently active. You can re-use it once it expires, or choose another.",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database Save Failed"})
	}

	// 5. Generate ShortCode via Base62/Hashids (Only if no custom code was provided)
	if newLink.ShortCode == "PENDING" {
		generatedCode := base62.Encode(uint64(newLink.LinkID))
		charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

		// SAFETY CHECK: Use a loop to guarantee we break any rare collisions
		for {
			var activeCount int64
			database.DB.Model(&models.Link{}).
				Where("short_code = ? AND expires_at > ? AND link_id != ?", generatedCode, time.Now().UTC(), newLink.LinkID).
				Count(&activeCount)

			// If count is 0, the code is completely free. Break the loop.
			if activeCount == 0 {
				break
			}

			// COLLISION! Append a random character from the charset to mutate it
			randomChar := string(charset[rand.Intn(len(charset))])
			generatedCode = generatedCode + randomChar
		}

		newLink.ShortCode = generatedCode

		// Persist the generated back to the DB
		if err := database.DB.Model(&newLink).Update("ShortCode", newLink.ShortCode).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update short code"})
		}
	}

	// 6. Redis Caching
	redisTTL := min(24*time.Hour, time.Until(expiryTime))
	_ = database.RedisDB.Set(ctx, newLink.ShortCode, newLink.LongURL, redisTTL).Err()

	// 7. Final Response
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
