package routes

import (
	"github.com/ankan8/url-shortener/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	// Create Redis client for database 0
	r := database.CreateClient(0)
	if r == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to connect to Redis"})
	}
	defer r.Close()

	// Fetch the value for the given URL key
	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short URL not found in database"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error connecting to the database"})
	}

	// Create Redis client for database 1
	rInr := database.CreateClient(1)
	if rInr == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to connect to Redis"})
	}
	defer rInr.Close()

	// Increment the counter for tracking hits
	if _, err := rInr.Incr(database.Ctx, "counter").Result(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to increment counter"})
	}

	// Redirect to the original URL
	return c.Redirect(value, fiber.StatusMovedPermanently)
}
