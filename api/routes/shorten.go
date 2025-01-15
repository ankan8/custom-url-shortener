package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/ankan8/url-shortener/database"
	"github.com/ankan8/url-shortener/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// Rate limiting
	r2 := database.CreateClient(1)
	if r2 == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to create Redis client"})
	}
	defer r2.Close()

	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		// Set rate limit quota if not set
		err = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*time.Minute).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to set rate limit quota"})
		}
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve rate limit data"})
	} else {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":            "rate limit exceeded",
				"rate_limit_reset": ttl / time.Minute,
			})
		}
	}

	// Validate URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid URL"})
	}

	// Remove domain errors
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid domain in URL"})
	}

	// Enforce HTTPS
	body.URL = helpers.EnforceHTTP(body.URL)

	// Generate short URL
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	// Check if the short URL already exists
	r := database.CreateClient(0)
	if r == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to create Redis client"})
	}
	defer r.Close()

	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "short URL already exists"})
	}

	// Set expiry if not provided
	if body.Expiry == 0 {
		body.Expiry = 24 * time.Hour
	}

	// Save the URL in Redis
	err = r.Set(database.Ctx, id, body.URL, body.Expiry).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to store URL"})
	}

	// Decrement rate limit
	if _, err := r2.Decr(database.Ctx, c.IP()).Result(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to decrement rate limit"})
	}

	// Prepare response
	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	xRateRemaining, _ := strconv.Atoi(val)
	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()

	resp := response{
		URL:             body.URL,
		CustomShort:     os.Getenv("DOMAIN") + "/" + id,
		Expiry:          body.Expiry,
		XRateRemaining:  xRateRemaining,
		XRateLimitReset: ttl / time.Minute,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
