package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zalg2261/bioskop/backend/internal/db"
	"github.com/zalg2261/bioskop/backend/internal/models"
)

func GetBranches(c *fiber.Ctx) error {
	city := c.Query("city")
	var branches []models.Branch
	query := db.DB
	
	if city != "" {
		query = query.Where("city = ?", city)
	}
	
	if err := query.Find(&branches).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(branches)
}

