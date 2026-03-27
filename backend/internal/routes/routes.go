package routes

import (
	"github.com/SKjustSK/alru-url-shortener/backend/internal/handlers"
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET(":short_code", handlers.RedirectLink)
	api := e.Group("/api")

	api.POST("/links", handlers.CreateLink)
	api.GET("/links", handlers.GetLinks)

	api.POST("/users", handlers.CreateUser)
	api.POST("/sessions", handlers.CreateSession)
}
