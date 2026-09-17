package main

import (
	"go-ecommerce/internal/config"
	"go-ecommerce/internal/database"
	"go-ecommerce/internal/handler"
	"go-ecommerce/internal/repository"
	"go-ecommerce/internal/router"
	"go-ecommerce/internal/service"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()
	addr := cfg.ServerAddress()

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	// 2. Connect DB
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	log.Println("Database connected")

	// 3. Wire dependencies (Repository -> Service -> Handler)
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg.JWTSecret, 24*time.Hour)
	userHandler := handler.NewUserHandler(userService)

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	// 4. Initialize Fiber App
	app := fiber.New()

	// 5. Setup Routes
	router.SetupRoutes(app, userHandler, productHandler, cfg.JWTSecret)

	// 6. Listen
	log.Printf("Server running on port %s", cfg.AppPort)
	if err := app.Listen(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
