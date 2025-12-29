package routes

import (
	"context"
	"strconv"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/runtimeninja/go-redis-url-shortner/database"
)

func Analytics(c *fiber.Ctx) error {
	short := c.Params("short")
	if short == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "short code is required",
		})
	}

	r := database.CreateClient(0)
	defer r.Close()

	ctx, cancel := context.WithTimeout(database.Ctx, 3*time.Second)
	defer cancel()

	clickKey := "clicks:" + short

	// read click count
	val, err := r.Get(ctx, clickKey).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"short":  short,
			"clicks": 0,
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read analytics data",
		})
	}

	clicks, err := strconv.Atoi(val)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "invalid analytics value",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"short":  short,
		"clicks": clicks,
	})
}
