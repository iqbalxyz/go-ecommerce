package main

import (
	"go-ecommerce/internal/config"
	"go-ecommerce/internal/database"
	"go-ecommerce/internal/models"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	db.AutoMigrate(&models.User{})

	app := fiber.New()

	log.Println("Database connected")

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

	log.Fatal(app.Listen(cfg.AppPort))
}
