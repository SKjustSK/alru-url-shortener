package links

import (
	"net/http"

	"github.com/SKjustSK/alru-url-shortener/backend/internal/database"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type TimeSeriesPoint struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// NEW: Struct for hourly drill-down data
type HourlyPoint struct {
	Date  string `json:"date"`
	Hour  string `json:"hour"`
	Count int64  `json:"count"`
}

type StatItem struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type AnalyticsResponse struct {
	ShortCode   string            `json:"short_code"`
	LongURL     string            `json:"long_url"`
	TotalClicks int64             `json:"total_clicks"`
	Timeline    []TimeSeriesPoint `json:"timeline"`
	Hourly      []HourlyPoint     `json:"hourly"`
	OS          []StatItem        `json:"os"`
	Browsers    []StatItem        `json:"browsers"`
	Devices     []StatItem        `json:"devices"`
	Countries   []StatItem        `json:"countries"`
	Referrers   []StatItem        `json:"referrers"`
}

func GetLinkAnalytics(c *echo.Context) error {
	shortCode := c.Param("short_code")

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*models.JWTCustomClaims)
	userID := claims.UserID

	var link models.Link
	if err := database.DB.Where("short_code = ? AND user_id = ?", shortCode, userID).First(&link).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Link not found or unauthorized"})
	}

	linkID := link.LinkID
	response := AnalyticsResponse{
		ShortCode: shortCode,
		LongURL:   link.LongURL,
	}

	database.DB.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&response.TotalClicks)

	if response.TotalClicks == 0 {
		response.Timeline = []TimeSeriesPoint{}
		response.Hourly = []HourlyPoint{}
		response.OS = []StatItem{}
		response.Browsers = []StatItem{}
		response.Devices = []StatItem{}
		response.Countries = []StatItem{}
		response.Referrers = []StatItem{}
		return c.JSON(http.StatusOK, response)
	}

	database.DB.Raw(`
		SELECT to_char(clicked_at, 'YYYY-MM-DD') as date, COUNT(*) as count 
		FROM clicks WHERE link_id = ? GROUP BY date ORDER BY date ASC
	`, linkID).Scan(&response.Timeline)

	database.DB.Raw(`
		SELECT to_char(clicked_at, 'YYYY-MM-DD') as date, to_char(clicked_at, 'HH24:00') as hour, COUNT(*) as count 
		FROM clicks WHERE link_id = ? GROUP BY date, hour ORDER BY date ASC, hour ASC
	`, linkID).Scan(&response.Hourly)

	database.DB.Raw(`SELECT os as name, COUNT(*) as count FROM clicks WHERE link_id = ? GROUP BY os ORDER BY count DESC`, linkID).Scan(&response.OS)
	database.DB.Raw(`SELECT browser as name, COUNT(*) as count FROM clicks WHERE link_id = ? GROUP BY browser ORDER BY count DESC`, linkID).Scan(&response.Browsers)
	database.DB.Raw(`SELECT device_type as name, COUNT(*) as count FROM clicks WHERE link_id = ? GROUP BY device_type ORDER BY count DESC`, linkID).Scan(&response.Devices)

	database.DB.Raw(`
		SELECT COALESCE(NULLIF(country, ''), 'Unknown') as name, COUNT(*) as count 
		FROM clicks WHERE link_id = ? GROUP BY name ORDER BY count DESC
	`, linkID).Scan(&response.Countries)

	database.DB.Raw(`
		SELECT COALESCE(NULLIF(referrer, ''), 'Direct / Email') as name, COUNT(*) as count 
		FROM clicks WHERE link_id = ? GROUP BY name ORDER BY count DESC
	`, linkID).Scan(&response.Referrers)

	return c.JSON(http.StatusOK, response)
}
