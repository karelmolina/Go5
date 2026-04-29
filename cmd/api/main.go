package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/karelmolina/play5/config"
	"github.com/karelmolina/play5/database"
)

func main() {
	database.ConnectDB()

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Play5",
	})

	app.Get("/health", func(c fiber.Ctx) error {
		sqlDB, err := database.DB.DB()
		if err != nil {
			return c.Status(503).JSON(fiber.Map{
				"status":   "degraded",
				"database": "disconnected",
			})
		}

		if err := sqlDB.Ping(); err != nil {
			return c.Status(503).JSON(fiber.Map{
				"status":   "degraded",
				"database": "disconnected",
			})
		}

		return c.JSON(fiber.Map{
			"status":   "ok",
			"database": "connected",
		})
	})

	port := config.Config("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":"+port, fiber.ListenConfig{EnablePrefork: true}))
}
