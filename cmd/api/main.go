package main

import (
	"go-ecommerce/internal/config"
	"go-ecommerce/internal/database"
	"go-ecommerce/internal/handler"
	"go-ecommerce/internal/repository"
	"go-ecommerce/internal/router"
	"go-ecommerce/internal/service"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()
	addr := cfg.ServerAddress()

	// 2. Connect DB
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	log.Println("Database connected")

	// 3. Wire dependencies (Repository -> Service -> Handler)
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// 4. Initialize Fiber App
	app := fiber.New()

	// 5. Setup Routes
	router.SetupRoutes(app, userHandler)

	// 6. Listen
	log.Printf("Server running on port %s", cfg.AppPort)
	if err := app.Listen(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
