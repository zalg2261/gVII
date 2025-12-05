package controllers

import (
    "time"
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
    "github.com/zalg2261/bioskop/backend/internal/db"
    "github.com/zalg2261/bioskop/backend/internal/models"
    "github.com/zalg2261/bioskop/backend/internal/services"
)

func RequestRefund(c *fiber.Ctx) error {
    id := c.Params("bookingId")
    var booking models.Booking
    if err := db.DB.First(&booking, id).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Booking not found"})
    }
    var body struct{ Reason string `json:"reason"` }
    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }
    if err := services.CreateRefund(booking, body.Reason); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"message": "Refund requested"})
}

func ApproveRefund(c *fiber.Ctx) error {
    id := c.Params("refundId")
    var refund models.Refund
    if err := db.DB.First(&refund, id).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Refund not found"})
    }
    if err := services.ApproveRefund(&refund); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"message": "Refund approved"})
}

func RefundBooking(c *fiber.Ctx) error {
	id := c.Params("bookingId")
	var booking models.Booking
	if err := db.DB.First(&booking, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Booking not found"})
	}

	if booking.Status != "PAID" {
		return c.Status(400).JSON(fiber.Map{"error": "Only PAID bookings can be refunded"})
	}

	var body struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		body.Reason = "Refund dari pihak bioskop"
	}

	var showtime models.Showtime
	db.DB.First(&showtime, booking.ShowtimeID)
	refundAmount := int64(booking.Seats) * showtime.Price

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// Update booking status
		booking.Status = "REFUNDED"
		if err := tx.Save(&booking).Error; err != nil {
			return err
		}

		// Release seats
		if err := tx.Model(&showtime).
			Update("seats_left", gorm.Expr("seats_left + ?", booking.Seats)).Error; err != nil {
			return err
		}

		// Create refund record
		refund := models.Refund{
			BookingID:   booking.ID,
			UserID:     booking.UserID,
			AmountCents: refundAmount,
			Reason:     body.Reason,
			Status:     "APPROVED",
			ProcessedAt: time.Now(),
		}
		if err := tx.Create(&refund).Error; err != nil {
			return err
		}

		// Update transaction status
		if err := tx.Model(&models.Transaction{}).
			Where("booking_id = ?", booking.ID).
			Update("status", "REFUND").Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Booking refunded successfully"})
}

func GetRefunds(c *fiber.Ctx) error {
	var refunds []models.Refund
	if err := db.DB.Preload("Booking").
		Order("created_at DESC").
		Find(&refunds).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(refunds)
}