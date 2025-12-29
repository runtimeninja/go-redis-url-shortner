package main

import (
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/runtimeninja/go-redis-url-shortner/routes"
)

func setupMiddlewares(app *fiber.App) {
	app.Use(logger.New())
}

func setupRoutes(app *fiber.App) {
	app.Post("/urlshortner/api/v1", routes.ShortenURL)
	app.Get("/urlshortner/analytics/v1/:short", routes.Analytics)
	app.Get("/:url", routes.ResolveURL)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	app := fiber.New()

	setupMiddlewares(app)

	setupRoutes(app)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("server running on PORT %s", port)

	log.Fatal(app.Listen(":" + port))
}