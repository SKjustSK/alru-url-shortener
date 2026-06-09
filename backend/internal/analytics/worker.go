package analytics

import (
	"context"
	"log"
	"time"

	"github.com/SKjustSK/alru-url-shortener/backend/internal/database"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/models"
)

var (
	clickChan chan models.Click
)

// Init initializes the buffered channel and starts the background worker.
func Init(ctx context.Context, bufferSize int, batchSize int, flushInterval time.Duration) {
	clickChan = make(chan models.Click, bufferSize)
	go worker(ctx, batchSize, flushInterval)
}

// QueueClick queues a click event. Uses non-blocking channel send to guarantee
// redirection performance even under extreme traffic spikes (load shedding).
func QueueClick(click models.Click) {
	if clickChan == nil {
		log.Println("Analytics worker not initialized. Dropping click event.")
		return
	}
	select {
	case clickChan <- click:
	default:
		log.Println("Analytics buffer full. Dropping click event to protect redirection latency.")
	}
}

func worker(ctx context.Context, batchSize int, flushInterval time.Duration) {
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	buffer := make([]models.Click, 0, batchSize)

	flush := func() {
		if len(buffer) == 0 {
			return
		}
		// Batch insert using GORM (GORM automatically converts slice inserts into a single bulk insert query)
		if err := database.DB.Create(&buffer).Error; err != nil {
			log.Printf("Failed to batch insert click analytics: %v", err)
		} else {
			log.Printf("Successfully flushed %d click records to PostgreSQL", len(buffer))
		}
		// Reset buffer slice while maintaining memory capacity
		buffer = buffer[:0]
	}

	for {
		select {
		case click, ok := <-clickChan:
			if !ok {
				flush()
				return
			}
			buffer = append(buffer, click)
			if len(buffer) >= batchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-ctx.Done():
			flush()
			return
		}
	}
}
