package worker

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/infrastructure/database"
	"log"
	"time"
)

func StartStatusExpiryWorker() {
	// Run cleanup every hour
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			cleanupExpiredStatuses()
		}
	}()
	log.Println("Status expiry worker started")
}

func cleanupExpiredStatuses() {
	now := time.Now()

	var expiredStatuses []entities.Status
	// Find all statuses that have expired
	err := database.DB.Where("expires_at <= ?", now).Find(&expiredStatuses).Error
	if err != nil {
		log.Printf("Status worker error fetching expired statuses: %v", err)
		return
	}

	if len(expiredStatuses) == 0 {
		log.Println("No expired statuses to clean up")
		return
	}

	// Delete expired statuses
	err = database.DB.Delete(&entities.Status{}, "expires_at <= ?", now).Error
	if err != nil {
		log.Printf("Status worker error deleting expired statuses: %v", err)
		return
	}

	log.Printf("Status worker cleaned up %d expired statuses", len(expiredStatuses))
}
