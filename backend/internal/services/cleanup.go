package services

import (
	"log"
	"time"

	"github.com/zalg2261/bioskop/backend/internal/db"
	"github.com/zalg2261/bioskop/backend/internal/models"
	"gorm.io/gorm"
)

func CleanupExpiredBookings() {
	log.Println("Starting cleanup of expired bookings...")

	var expiredBookings []models.Booking
	if err := db.DB.Where("status = ? AND expires_at < ?", "PENDING", time.Now()).Find(&expiredBookings).Error; err != nil {
		log.Printf("Error finding expired bookings: %v", err)
		return
	}

	if len(expiredBookings) == 0 {
		log.Println("No expired bookings to cleanup")
		return
	}

	log.Printf("Found %d expired bookings to cleanup", len(expiredBookings))

	for _, booking := range expiredBookings {
		err := db.DB.Transaction(func(tx *gorm.DB) error {
			// Update booking status to CANCELLED
			if err := tx.Model(&booking).Update("status", "CANCELLED").Error; err != nil {
				return err
			}

			// Release seats
			if err := tx.Model(&models.Showtime{}).
				Where("id = ?", booking.ShowtimeID).
				Update("seats_left", gorm.Expr("seats_left + ?", booking.Seats)).Error; err != nil {
				return err
			}

			// Delete seat lock
			if err := tx.Where("booking_id = ?", booking.ID).Delete(&models.SeatLock{}).Error; err != nil {
				return err
			}

			log.Printf("Cleaned up booking ID: %d, released %d seats", booking.ID, booking.Seats)
			return nil
		})

		if err != nil {
			log.Printf("Error cleaning up booking ID %d: %v", booking.ID, err)
		}
	}

	log.Println("Cleanup completed")
}

func StartCleanupJob() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			CleanupExpiredBookings()
		}
	}()
	log.Println("Cleanup job started")
}

