package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/zalg2261/bioskop/backend/internal/db"
	"github.com/zalg2261/bioskop/backend/internal/models"
)

// POST /wallet/topup
// body: { "amount": 100000 } // in rupiah
func TopUpWallet(c *fiber.Ctx) error {
	// Get user_id from JWT token
	userIDStr := c.Locals("user_id").(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)

	var body struct {
		Amount int64 `json:"amount"` // in rupiah
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if body.Amount <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Amount must be greater than 0"})
	}

	// Update user balance
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	user.Balance += body.Amount
	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Topup success", "new_balance": user.Balance})
}
