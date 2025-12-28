package routes

import (
	"context"
    "log"
    "time"
    "github.com/runtimeninja/go-redis-url-shortner/database"
    "github.com/redis/go-redis/v9"
    "github.com/gofiber/fiber/v2"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	r := database.CreateClient(0)
	defer r.Close()

	// local timeout context for Redis operations
	timeoutCtx, cancel := context.WithTimeout(database.Ctx, 3*time.Second)
	defer cancel()

	value, err := r.Get(timeoutCtx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "short not found in database",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot connect to database",
		})
	}

	if err := r.Incr(timeoutCtx, "counter").Err(); err != nil {
		log.Printf("failed to increment counter: %v", err)
	}

	return c.Redirect(value, 301)
}
