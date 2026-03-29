package links

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func GetLinks(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "This route is yet to be implemented."})
}
