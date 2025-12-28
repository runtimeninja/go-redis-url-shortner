package routes

import (
	"context"
    "log"
    "os"
    "strconv"
    "time"
    "github.com/runtimeninja/go-redis-url-shortner/helpers"
    "github.com/runtimeninja/go-redis-url-shortner/database"
    "github.com/redis/go-redis/v9"
    "github.com/gofiber/fiber/v2"
    "github.com/asaskevich/govalidator"
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
	XRateLimitReset int           `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// rate limit db
	rateClient := database.CreateClient(1)
	defer rateClient.Close()

	ip := c.IP()
	timeoutCtx, cancel := context.WithTimeout(database.Ctx, 3*time.Second)
	defer cancel()

	val, err := rateClient.Get(timeoutCtx, ip).Result()
	if err == redis.Nil {
		apiQuota := os.Getenv("API_QUOTA")
		if err := rateClient.Set(timeoutCtx, ip, apiQuota, 30*time.Minute).Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot initialize quota"})
		}
		val = apiQuota
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch rate limit"})
	}

	quota, err := strconv.Atoi(val)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid rate limit value"})
	}

	if quota <= 0 {
		ttl, _ := rateClient.TTL(timeoutCtx, ip).Result()
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":            "Rate Limit exceeded",
			"rate_limit_reset": int(ttl.Minutes()),
		})
	}

	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "You can't access this domain"})
	}

	body.URL = helpers.EnforceHTTP(body.URL)

	// generate short id
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	// main url db
	urlClient := database.CreateClient(0)
	defer urlClient.Close()

	existing, err := urlClient.Get(timeoutCtx, id).Result()
	if err != nil && err != redis.Nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database error"})
	}
	if existing != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Custom short already in use"}) // short ID already exists
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	if err := urlClient.Set(timeoutCtx, id, body.URL, body.Expiry*time.Hour).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to save URL"})
	}

	if err := rateClient.Decr(timeoutCtx, ip).Err(); err != nil {
		log.Println("Failed to decrement quota", err)
	}

	newVal, err := rateClient.Get(timeoutCtx, ip).Result()
	if err != nil {
		newVal = "0"
	}

	rateRemaining, _ := strconv.Atoi(newVal)
	ttl, _ := rateClient.TTL(timeoutCtx, ip).Result()

	resp := response{
		URL:             body.URL,
		CustomShort:     os.Getenv("DOMAIN") + "/" + id,
		Expiry:          body.Expiry,
		XRateRemaining:  rateRemaining,
		XRateLimitReset: int(ttl.Minutes()),
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
