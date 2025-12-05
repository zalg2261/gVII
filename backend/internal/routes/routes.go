package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zalg2261/bioskop/backend/internal/controllers"
	"github.com/zalg2261/bioskop/backend/internal/middleware"
)

func SetupRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "gVII API",
			"version": "1.0.0",
		})
	})

	// PUBLIC ROUTES
	// AUTH
	auth := app.Group("/auth")
	auth.Post("/login", controllers.Login)
	auth.Post("/register", controllers.Register)

	// PUBLIC: Browse movies and schedules (no auth required)
	app.Get("/movies", controllers.GetMovies)
	app.Get("/movies/:id", controllers.GetMovie)
	app.Get("/schedule", controllers.GetSchedules) // Public: browse schedules
	app.Get("/schedule/:id", controllers.GetSchedule) // Public: get single schedule
	app.Get("/branches", controllers.GetBranches)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// PROTECTED ROUTES (require authentication)
	// BOOKING (user only)
	app.Post("/book", middleware.RequireAuth, controllers.CreateBooking)
	app.Post("/payment/:bookingId", middleware.RequireAuth, controllers.CompletePayment)
	app.Post("/payment/failed/:bookingId", middleware.RequireAuth, controllers.PaymentFailed)
	app.Get("/my-bookings", middleware.RequireAuth, controllers.GetMyBookings)
	app.Post("/wallet/topup", middleware.RequireAuth, controllers.TopUpWallet)

	// ADMIN ROUTES (require admin role)
	admin := app.Group("/admin", middleware.RequireAdmin)
	
	// CRUD Jadwal Tayang (Admin only)
	admin.Post("/schedule", controllers.CreateSchedule)
	admin.Put("/schedule/:id", controllers.UpdateSchedule)
	admin.Delete("/schedule/:id", controllers.DeleteSchedule)
	admin.Post("/schedule/:id/cancel", controllers.CancelShowtime) // Cancel showtime and refund
	
	// CRUD Movies (Admin only)
	admin.Post("/movies", controllers.CreateMovie)
	admin.Put("/movies/:id", controllers.UpdateMovie)
	admin.Delete("/movies/:id", controllers.DeleteMovie)
	
	// Refund management (Admin only)
	admin.Post("/refund/:bookingId", controllers.RefundBooking)
	admin.Get("/refunds", controllers.GetRefunds)

}
