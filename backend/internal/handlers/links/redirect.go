package links

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"time"

	"github.com/SKjustSK/alru-url-shortener/backend/internal/database"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/models"
	"github.com/labstack/echo/v5"
	"github.com/mileusna/useragent"
)

// 1. The Route Handlers
func RedirectGenerated(c *echo.Context) error {
	shortCode := c.Param("short_code")
	return processRedirect(c, shortCode, false)
}

func RedirectCustom(c *echo.Context) error {
	alias := c.Param("alias")
	return processRedirect(c, alias, true)
}

// 2. The Core Redirect Logic
func processRedirect(c *echo.Context, code string, isCustom bool) error {
	ctx := c.Request().Context()

	ip := c.RealIP()
	uaString := c.Request().UserAgent()
	referrer := c.Request().Referer()
	country := c.Request().Header.Get("CF-IPCountry")

	var redirectURL string
	var linkID int64

	// Namespace the Redis key based on the link type
	redisKey := "gen:" + code
	if isCustom {
		redisKey = "custom:" + code
	}

	// Check Redis
	longURL, err := database.RedisDB.Get(ctx, redisKey).Result()
	if err == nil {
		// Cache Hit
		redirectURL = longURL
	} else {
		// Cache Miss - Query DB using BOTH code and IsCustom flag
		var link models.Link
		err = database.DB.Select("link_id", "long_url", "expires_at").
			Where("short_code = ? AND is_custom = ? AND expires_at > ?", code, isCustom, time.Now().UTC()).
			Order("expires_at DESC").
			First(&link).Error

		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Link not found or expired"})
		}

		redirectURL = link.LongURL
		linkID = link.LinkID

		// Re-populate Redis
		redisTTL := min(24*time.Hour, time.Until(link.ExpiresAt))
		_ = database.RedisDB.Set(ctx, redisKey, link.LongURL, redisTTL).Err()
	}

	// Background analytics worker (passing isCustom so it finds the correct ID on a cache hit)
	go trackClick(code, isCustom, linkID, ip, uaString, referrer, country)

	// Redirect
	return c.Redirect(http.StatusFound, redirectURL)
}

// 3. The Analytics Worker
func trackClick(shortCode string, isCustom bool, linkID int64, ip, uaString, referrer, country string) {
	bgCtx := context.Background()

	// 1. Get LinkID if it was a Cache Hit
	if linkID == 0 {
		var link models.Link
		// Crucial Fix: Must include is_custom here, otherwise it might fetch the wrong ID!
		if err := database.DB.WithContext(bgCtx).Select("link_id").
			Where("short_code = ? AND is_custom = ?", shortCode, isCustom).First(&link).Error; err != nil {
			return
		}
		linkID = link.LinkID
	}

	// 2. Hash the IP for Privacy
	salt := os.Getenv("IP_HASH_SALT")
	if salt == "" {
		salt = "fallback-alru-salt"
	}
	hash := sha256.Sum256([]byte(ip + salt))
	ipHash := hex.EncodeToString(hash[:])

	// 3. Parse the User Agent cleanly
	ua := useragent.Parse(uaString)

	var deviceType string
	switch {
	case ua.Bot:
		deviceType = "Bot"
	case ua.Tablet:
		deviceType = "Tablet"
	case ua.Mobile:
		deviceType = "Mobile"
	default:
		deviceType = "Desktop"
	}

	// 4. Build the Click Record
	click := models.Click{
		LinkID:     linkID,
		IPHash:     ipHash,
		Country:    country,
		Referrer:   referrer,
		UserAgent:  uaString,
		DeviceType: deviceType,
		OS:         ua.OS,
		Browser:    ua.Name,
	}

	// 5. Save to Postgres
	database.DB.WithContext(bgCtx).Create(&click)
}
