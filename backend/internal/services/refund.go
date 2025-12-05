package services

import (
    "errors"
    "time"

    "github.com/zalg2261/bioskop/backend/internal/db"
    "github.com/zalg2261/bioskop/backend/internal/models"
    "gorm.io/gorm"
)


// ======================================
// Request Refund (dipanggil dari controller)
// ======================================
func CreateRefund(booking models.Booking, reason string) error {
    // hanya booking PAID yang bisa refund
    if booking.Status != "PAID" {
        return errors.New("booking is not paid")
    }

    refund := models.Refund{
        BookingID: booking.ID,
        Reason:    reason,
        Status:    "PENDING",
        CreatedAt: time.Now(),
    }

    return db.DB.Create(&refund).Error
}



func ApproveRefund(refund *models.Refund) error {
	var booking models.Booking
	if err := db.DB.First(&booking, refund.BookingID).Error; err != nil {
		return err
	}

	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Update refund status
		refund.Status = "APPROVED"
		refund.ProcessedAt = time.Now()
		if err := tx.Save(refund).Error; err != nil {
			return err
		}

		// Update booking status
		booking.Status = "REFUNDED"
		if err := tx.Save(&booking).Error; err != nil {
			return err
		}

		// Add balance back to user
		var user models.User
		if err := tx.First(&user, refund.UserID).Error; err != nil {
			return err
		}
		user.Balance += refund.AmountCents
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// Release seats
		var showtime models.Showtime
		if err := tx.First(&showtime, booking.ShowtimeID).Error; err != nil {
			return err
		}
		showtime.SeatsLeft += booking.Seats
		if err := tx.Save(&showtime).Error; err != nil {
			return err
		}

		return nil
	})
}
