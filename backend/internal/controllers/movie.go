package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zalg2261/bioskop/backend/internal/db"
	"github.com/zalg2261/bioskop/backend/internal/models"
)

func GetMovies(c *fiber.Ctx) error {
	var movies []models.Movie
	if err := db.DB.Find(&movies).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(movies)
}

func GetMovie(c *fiber.Ctx) error {
	id := c.Params("id")
	var movie models.Movie
	if err := db.DB.First(&movie, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Movie not found"})
	}
	return c.JSON(movie)
}

func CreateMovie(c *fiber.Ctx) error {
	var movie models.Movie
	if err := c.BodyParser(&movie); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := db.DB.Create(&movie).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(movie)
}

func UpdateMovie(c *fiber.Ctx) error {
	id := c.Params("id")
	var movie models.Movie
	if err := db.DB.First(&movie, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Movie not found"})
	}
	var updateData models.Movie
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := db.DB.Model(&movie).Updates(updateData).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(movie)
}

func DeleteMovie(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.DB.Delete(&models.Movie{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Movie deleted"})
}

