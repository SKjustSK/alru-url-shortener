package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/SKjustSK/alru-url-shortener/backend/internal/database"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("internal/config/.env"); err != nil {
		log.Println("No .env file found. Relying on system environment variables.")
	}

	// Connect databases
	database.ConnectPostgreSQL()
	database.ConnectRedis()

	e := echo.New()

	// Global Middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			os.Getenv("FRONTEND_URL"),
			"http://localhost:5173",
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowCredentials: true,
	}))

	// Register all our routes
	routes.Register(e)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}

	sc := echo.StartConfig{Address: ":" + port}
	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
