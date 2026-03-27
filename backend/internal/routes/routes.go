package routes

import (
	"github.com/SKjustSK/alru-url-shortener/backend/internal/handlers/links"
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Echo) {
	// e.GET(":short_code", handlers.RedirectLink)
	api := e.Group("/api")

	api.POST("/links", links.CreateLink)
	// api.GET("/links", handlers.GetLinks)

	// api.POST("/users", handlers.CreateUser)
	// api.POST("/sessions", handlers.CreateSession)
}
