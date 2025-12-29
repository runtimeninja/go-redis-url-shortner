package routes

import (
	"context"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/runtimeninja/go-redis-url-shortner/database"
)

func ResolveURL(c *fiber.Ctx) error {
	short := c.Params("url")

	r := database.CreateClient(0)
	defer r.Close()

	ctx, cancel := context.WithTimeout(database.Ctx, 3*time.Second)
	defer cancel()

	originalURL, err := r.Get(ctx, short).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "short url not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "database error",
		})
	}

	_ = r.Incr(ctx, "clicks:"+short).Err()

	return c.Redirect(originalURL, fiber.StatusMovedPermanently)
}
