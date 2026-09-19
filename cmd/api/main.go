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
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	fiberRecover "github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	addr := cfg.ServerAddress()

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

	cartRepo := repository.NewCartRepository(db)
	cartService := service.NewCartService(cartRepo, productRepo)
	cartHandler := handler.NewCartHandler(cartService)

	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, cartRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	// 4. Initialize Fiber App
	app := fiber.New()

	app.Use(logger.New())
	app.Use(fiberRecover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: false,
		MaxAge:           3600,
	}))

	router.SetupRoutes(
		app,
		userHandler,
		productHandler,
		cartHandler,
		orderHandler,
		cfg.JWTSecret,
	)

	// 6. Listen
	log.Printf("Server running on port %s", cfg.AppPort)
	if err := app.Listen(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
