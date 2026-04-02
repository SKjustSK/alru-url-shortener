package links

import (
	"net/http"
	"os"
	"time"

	"github.com/SKjustSK/alru-url-shortener/backend/internal/database"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type LinkItem struct {
	ShortURL  string    `json:"short_url"`
	LongURL   string    `json:"long_url"`
	ShortCode string    `json:"short_code"`
	IsCustom  bool      `json:"is_custom"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type GetLinksResponse struct {
	Links []LinkItem `json:"links"`
	Total int        `json:"total"`
}

func GetLinks(c *echo.Context) error {
	// 1. Extract the authenticated UserID from the JWT context
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*models.JWTCustomClaims)
	userID := claims.UserID

	// 2. Query the database
	var userLinks []models.Link

	if err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&userLinks).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve links from the database",
		})
	}

	// 3. Format the response
	baseURL := os.Getenv("BACKEND_URL")
	if baseURL == "" {
		baseURL = "http://localhost:1323"
	}

	formattedLinks := make([]LinkItem, 0, len(userLinks))

	for _, link := range userLinks {
		// Determine the correct route prefix
		prefix := "/"
		if link.IsCustom {
			prefix = "/c/"
		}

		formattedLinks = append(formattedLinks, LinkItem{
			ShortURL:  baseURL + prefix + link.ShortCode,
			LongURL:   link.LongURL,
			ShortCode: link.ShortCode,
			IsCustom:  link.IsCustom,
			CreatedAt: link.CreatedAt,
			ExpiresAt: link.ExpiresAt,
		})
	}

	// 4. Return the data
	return c.JSON(http.StatusOK, GetLinksResponse{
		Links: formattedLinks,
		Total: len(formattedLinks),
	})
}
