package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/rohitchauraisa1997/url-shortener/database"
	"github.com/rohitchauraisa1997/url-shortener/routes"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", routes.IsAlive)
	app.Get("/admin/route/resolutions/analytics", routes.GetUrlResolutionAnalytics)
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/shorten", routes.ShortenURL)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	database.CreateClient(0)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	setupRoutes(app)
	app.Listen(os.Getenv("APP_PORT"))
}
