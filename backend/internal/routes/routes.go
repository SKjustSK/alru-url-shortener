package routes

import (
	"net/http"
	"os"

	"github.com/SKjustSK/alru-url-shortener/backend/internal/handlers/auth"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/handlers/links"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/middleware"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Echo) {
	e.GET("/ping", func(c *echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	api := e.Group("/api")
	{
		api.POST("/users", auth.CreateUser)
		api.POST("/sessions", auth.CreateSession)

		// JWT Middleware
		jwtConfig := echojwt.Config{
			SigningKey: []byte(os.Getenv("JWT_SECRET")),
			NewClaimsFunc: func(c *echo.Context) jwt.Claims {
				return new(models.JWTCustomClaims)
			},
		}
		protected := api.Group("", echojwt.WithConfig(jwtConfig))
		{
			protected.POST("/links", links.CreateLink, middleware.LinkRateLimiter)
			protected.GET("/links", links.GetLinks)
			protected.GET("/links/:short_code/analytics", links.GetLinkAnalytics)
		}
	}

	e.GET("/:short_code", links.RedirectLink)
}
