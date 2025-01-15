package main

import (
	
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/ankan8/url-shortener/routes"
)

func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v\n", err)
	}

	// Initialize Fiber app
	app := fiber.New()

	// Use logger middleware
	app.Use(logger.New())

	// Set up routes
	setupRoutes(app)

	// Get port from environment or use default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = ":3000" // Default to port 3000 if not set
	}

	// Start the server
	log.Printf("Starting server on %s\n", port)
	log.Fatal(app.Listen(port))
}
