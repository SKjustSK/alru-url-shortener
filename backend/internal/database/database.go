package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/SKjustSK/alru-url-shortener/backend/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisDB *redis.Client

func ConnectPostgreSQL() *gorm.DB {
	var dsn string

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL != "" {
		dsn = databaseURL
		log.Println("Using provided DATABASE_URL.")
	} else {
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		sslmode := os.Getenv("DB_SSL")

		if host == "" || user == "" || dbname == "" || port == "" {
			log.Fatal("Database Error: Missing PostgreSQL connection variables")
		}

		// Default to disable if not explicitly set
		if sslmode == "" {
			sslmode = "disable"
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			host, user, password, dbname, port, sslmode,
		)
		log.Println("Constructed PostgreSQL DSN internally.")
	}

	// Connect PostgreSQL database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Auto migrate tables to PostgreSQL
	db.AutoMigrate(&models.Link{}, &models.User{})

	log.Println("PostgreSQL connection established.")
	DB := db
}

func ConnectRedis() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")

	if redisURL == "" {
		host := os.Getenv("REDIS_HOST")
		port := os.Getenv("REDIS_PORT")
		password := os.Getenv("REDIS_PASSWORD")

		if host == "" || port == "" {
			log.Fatal("Database Error: Missing Redis connection variables.")
		}

		redisURL = fmt.Sprintf("redis://:%s@%s:%s/0", password, host, port)

		log.Printf("Constructed Redis URL internally.")
	} else {
		log.Printf("Using provided Redis URL.")
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	rdb := redis.NewClient(opts)

	// Test connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connection established.")
	RedisDB := rdb
}
