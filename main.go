package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize the database
	db := InitDB()

	// Create the Fiber app instance
	app := fiber.New(fiber.Config{
		AppName: "Library API",
	})

	// Define the Auth routes. This will be public
	AuthHandlers(app.Group("/auth"), db)

	// Verify the JWT token. If valid, it will set the user id in the context
	protected := app.Use(AuthMiddleware(db))

	// Define the Book routes. This will be protected. Requires a valid JWT
	BookHandlers(protected.Group("/book"), db)

	// Download the books as either CSV or JSON
	Download(protected.Group("/download"), db)

	// Start the server on port 3000
	app.Listen(":3000")
}
