package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/SKjustSK/alru-url-shortener/backend/internal/analytics"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/database"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/routes"
)

func main() {
	// Set up cancelable context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals in a separate goroutine
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutdown signal received. Shutting down gracefully...")
		cancel()
	}()

	// Load environment variables
	if err := godotenv.Load("internal/config/.env"); err != nil {
		log.Println("No .env file found. Relying on system environment variables.")
	}

	// Connect databases
	database.ConnectPostgreSQL()
	database.ConnectRedis()

	// Initialize the asynchronous write-behind analytics worker
	// Buffer size of 5000, batch size of 100, flush interval of 5 seconds
	analytics.Init(ctx, 5000, 100, 5*time.Second)

	e := echo.New()

	// Global Middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// CORS middleware
	frontendURL := os.Getenv("FRONTEND_URL")
	allowedOrigins := []string{
		"http://localhost",
		"http://localhost:80",
		"http://localhost:5173",
	}
	if frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowedOrigins,
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
	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
