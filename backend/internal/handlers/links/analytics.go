package links

import (
	"net/http"
	"os"

	"github.com/SKjustSK/alru-url-shortener/backend/internal/database"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

// TimeSeriesPoint represents daily click counts for line charts
type TimeSeriesPoint struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// StatItem represents a generic grouped stat (e.g., "Chrome": 50, "Mobile": 120) for pie charts
type StatItem struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// AnalyticsResponse contains all the aggregated data for the React dashboard
type AnalyticsResponse struct {
	ShortURL    string            `json:"short_code"`
	LongURL     string            `json:"long_url"`
	TotalClicks int64             `json:"total_clicks"`
	Timeline    []TimeSeriesPoint `json:"timeline"`
	OS          []StatItem        `json:"os"`
	Browsers    []StatItem        `json:"browsers"`
	Devices     []StatItem        `json:"devices"`
	Countries   []StatItem        `json:"countries"`
	Referrers   []StatItem        `json:"referrers"`
}

func GetLinkAnalytics(c *echo.Context) error {
	shortCode := c.Param("short_code")

	// 1. Extract the user from the JWT
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*models.JWTCustomClaims)
	userID := claims.UserID

	// 2. Security Check: Does this user actually own this link?
	var link models.Link
	if err := database.DB.Where("short_code = ? AND user_id = ?", shortCode, userID).First(&link).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Link not found or unauthorized"})
	}

	linkID := link.LinkID
	response := AnalyticsResponse{
		ShortURL: os.Getenv("BACKEND_URL") + "/" + shortCode,
		LongURL:  link.LongURL,
	}

	// 3. Get Total Clicks
	database.DB.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&response.TotalClicks)

	// If there are no clicks yet, return early with empty arrays so the frontend doesn't crash
	if response.TotalClicks == 0 {
		response.Timeline = []TimeSeriesPoint{}
		response.OS = []StatItem{}
		response.Browsers = []StatItem{}
		response.Devices = []StatItem{}
		response.Countries = []StatItem{}
		response.Referrers = []StatItem{}
		return c.JSON(http.StatusOK, response)
	}

	// 4. Time-Series Data (Line Chart)
	// We use PostgreSQL's to_char to format the timestamp cleanly into "YYYY-MM-DD"
	database.DB.Raw(`
		SELECT to_char(clicked_at, 'YYYY-MM-DD') as date, COUNT(*) as count 
		FROM clicks 
		WHERE link_id = ? 
		GROUP BY date 
		ORDER BY date ASC
	`, linkID).Scan(&response.Timeline)

	// 5. Categorical Data (Pie / Bar Charts)
	// We run these queries to group the data and sort by the most popular first

	database.DB.Raw(`SELECT os as name, COUNT(*) as count FROM clicks WHERE link_id = ? GROUP BY os ORDER BY count DESC`, linkID).Scan(&response.OS)
	database.DB.Raw(`SELECT browser as name, COUNT(*) as count FROM clicks WHERE link_id = ? GROUP BY browser ORDER BY count DESC`, linkID).Scan(&response.Browsers)
	database.DB.Raw(`SELECT device_type as name, COUNT(*) as count FROM clicks WHERE link_id = ? GROUP BY device_type ORDER BY count DESC`, linkID).Scan(&response.Devices)

	// For Referrers and Countries, empty strings usually mean "Direct Click" or "Unknown IP"
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
