package controllers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/zalg2261/bioskop/backend/internal/db"
	"github.com/zalg2261/bioskop/backend/internal/models"
)


func CreateBooking(c *fiber.Ctx) error {
	// Get user_id from JWT token
	userIDStr := c.Locals("user_id").(string)
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user"})
	}

	var body struct {
		ShowtimeID uint `json:"showtime_id"`
		Seats      int  `json:"seats"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if body.Seats <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Seats must be greater than 0"})
	}

	// Check showtime availability
	var showtime models.Showtime
	if err := db.DB.First(&showtime, body.ShowtimeID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Showtime not found"})
	}

	if showtime.Status != "ACTIVE" {
		return c.Status(400).JSON(fiber.Map{"error": "Showtime is not active"})
	}

	if showtime.SeatsLeft < body.Seats {
		return c.Status(400).JSON(fiber.Map{"error": "Not enough seats available"})
	}

	expires := time.Now().Add(10 * time.Minute)

	booking := models.Booking{
		UserID:     uint(userID),
		ShowtimeID: body.ShowtimeID,
		Seats:      body.Seats,
		Status:     "PENDING",
		ExpiresAt:  &expires,
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		// Reduce seats_left
		if err := tx.Model(&showtime).Update("seats_left", gorm.Expr("seats_left - ?", body.Seats)).Error; err != nil {
			return err
		}

		// Create booking
		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		// Create seat lock
		lock := models.SeatLock{
			ShowtimeID: body.ShowtimeID,
			BookingID:  &booking.ID,
			SeatsCount: body.Seats,
			LockedAt:   time.Now(),
			ExpiresAt:  expires,
		}

		if err := tx.Create(&lock).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Reload booking with relations
	db.DB.Preload("Showtime").Preload("Showtime.Movie").First(&booking, booking.ID)

	return c.Status(fiber.StatusCreated).JSON(booking)
}



func CompletePayment(c *fiber.Ctx) error {
	// Get user_id from JWT token
	userIDStr := c.Locals("user_id").(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)

	id := c.Params("bookingId")
	bookingID, _ := strconv.Atoi(id)

	var booking models.Booking
	if err := db.DB.First(&booking, bookingID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Booking not found"})
	}

	// Check if booking belongs to user
	if booking.UserID != uint(userID) {
		return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Check if booking is still pending and not expired
	if booking.Status != "PENDING" {
		return c.Status(400).JSON(fiber.Map{"error": "Booking is not pending"})
	}

	if booking.ExpiresAt != nil && booking.ExpiresAt.Before(time.Now()) {
		return c.Status(400).JSON(fiber.Map{"error": "Booking has expired"})
	}

	// Get showtime for price
	var showtime models.Showtime
	db.DB.First(&showtime, booking.ShowtimeID)
	totalAmount := int64(booking.Seats) * showtime.Price

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// Update booking status
		booking.Status = "PAID"
		if err := tx.Save(&booking).Error; err != nil {
			return err
		}

		// Create transaction record
		txRecord := models.Transaction{
			BookingID: uint(bookingID),
			Amount:    totalAmount,
			Status:    "SUCCESS",
		}
		if err := tx.Create(&txRecord).Error; err != nil {
			return err
		}

		// Delete seat lock (seats are now permanently locked)
		if err := tx.Where("booking_id = ?", bookingID).Delete(&models.SeatLock{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Payment completed", "booking": booking})
}



func PaymentFailed(c *fiber.Ctx) error {
	id := c.Params("bookingId")
	bookingID, _ := strconv.Atoi(id)

	var booking models.Booking
	if err := db.DB.First(&booking, bookingID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Booking not found"})
	}

	if booking.Status != "PENDING" {
		return c.Status(400).JSON(fiber.Map{"error": "Booking is not pending"})
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// Release seats
		if err := tx.Model(&models.Showtime{}).
			Where("id = ?", booking.ShowtimeID).
			Update("seats_left", gorm.Expr("seats_left + ?", booking.Seats)).Error; err != nil {
			return err
		}

		// Update booking status
		booking.Status = "CANCELLED"
		if err := tx.Save(&booking).Error; err != nil {
			return err
		}

		// Delete seat lock
		if err := tx.Where("booking_id = ?", bookingID).Delete(&models.SeatLock{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Payment failed & seats released"})
}

func GetMyBookings(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)

	var bookings []models.Booking
	if err := db.DB.Where("user_id = ?", userID).
		Preload("Showtime").
		Preload("Showtime.Movie").
		Preload("Showtime.Branch").
		Order("created_at DESC").
		Find(&bookings).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(bookings)
}
