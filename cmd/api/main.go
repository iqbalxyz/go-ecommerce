package main

import (
	"go-ecommerce/internal/config"
	"go-ecommerce/internal/database"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg := config.LoadConfig()

	_, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	app := fiber.New()

	log.Println("Database connected")

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

	log.Fatal(app.Listen(cfg.AppPort))
}
