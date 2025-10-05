package main

import (
	"log"
	"os"

	"github.com/Blaze5333/shorten-url/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setUpRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}
	app := fiber.New()
	app.Use(logger.New())
	setUpRoutes(app)
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
