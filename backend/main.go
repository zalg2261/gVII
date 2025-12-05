package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/zalg2261/bioskop/backend/internal/db"
	"github.com/zalg2261/bioskop/backend/internal/routes"
	"github.com/zalg2261/bioskop/backend/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not loaded, relying on environment")
	}

	// connect DB
	db.Connect()

	// Start background cleanup job
	services.StartCleanupJob()

	app := fiber.New()

	// Enable CORS
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		
		// Handle preflight
		if c.Method() == "OPTIONS" {
			return c.SendStatus(200)
		}
		
		return c.Next()
	})

	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	log.Println("Server running on :" + port)
	log.Fatal(app.Listen(":" + port))
}
