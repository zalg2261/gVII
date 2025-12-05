package controllers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/zalg2261/bioskop/backend/internal/db"
    "github.com/zalg2261/bioskop/backend/internal/models"
    "gorm.io/gorm"
    "time"
)

func GetSchedules(c *fiber.Ctx) error {
    city := c.Query("city")
    branchID := c.Query("branch_id")
    movieID := c.Query("movie_id")

    var schedules []models.Showtime
    query := db.DB.Preload("Movie").Preload("Branch").Where("showtimes.status = ?", "ACTIVE")

    if city != "" {
        query = query.Joins("JOIN branches ON branches.id = showtimes.branch_id").
            Where("branches.city = ?", city)
    }
    if branchID != "" {
        query = query.Where("branch_id = ?", branchID)
    }
    if movieID != "" {
        query = query.Where("movie_id = ?", movieID)
    }

    query = query.Where("show_time > ?", time.Now())

    if err := query.Order("show_time ASC").Find(&schedules).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(schedules)
}

func GetSchedule(c *fiber.Ctx) error {
    id := c.Params("id")
    var schedule models.Showtime
    if err := db.DB.Preload("Movie").First(&schedule, id).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Schedule not found"})
    }
    return c.JSON(schedule)
}

func CreateSchedule(c *fiber.Ctx) error {
    var schedule models.Showtime
    if err := c.BodyParser(&schedule); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }
    if schedule.SeatsLeft == 0 && schedule.SeatsTotal > 0 {
        schedule.SeatsLeft = schedule.SeatsTotal
    }
    if err := db.DB.Create(&schedule).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.Status(201).JSON(schedule)
}

func UpdateSchedule(c *fiber.Ctx) error {
    id := c.Params("id")
    var schedule models.Showtime
    if err := db.DB.First(&schedule, id).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Schedule not found"})
    }
    var updateData models.Showtime
    if err := c.BodyParser(&updateData); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }
    if err := db.DB.Model(&schedule).Updates(updateData).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(schedule)
}

func DeleteSchedule(c *fiber.Ctx) error {
    id := c.Params("id")
    if err := db.DB.Delete(&models.Showtime{}, id).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"message": "Deleted"})
}

func CancelShowtime(c *fiber.Ctx) error {
	id := c.Params("id")
	var showtime models.Showtime
	if err := db.DB.First(&showtime, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Showtime not found"})
	}

	var body struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		body.Reason = "Kendala teknis dari pihak bioskop"
	}

	var bookings []models.Booking
	if err := db.DB.Where("showtime_id = ? AND status = ?", id, "PAID").Find(&bookings).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// Update showtime status
		showtime.Status = "CANCELLED"
		if err := tx.Save(&showtime).Error; err != nil {
			return err
		}

		// Refund all bookings
		for _, booking := range bookings {
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
			var showtimeForPrice models.Showtime
			tx.First(&showtimeForPrice, booking.ShowtimeID)
			refundAmount := int64(booking.Seats) * showtimeForPrice.Price

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
		}

		return nil
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Showtime cancelled and refunds processed",
		"refunded_bookings": len(bookings),
	})
}