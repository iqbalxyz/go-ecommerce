package router

import (
	"go-ecommerce/internal/handler"
	"go-ecommerce/internal/middleware"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func SetupRoutes(
	app *fiber.App,
	userHandler *handler.UserHandler,
	productHandler *handler.ProductHandler,
	cartHandler *handler.CartHandler,
	orderHandler *handler.OrderHandler,
	jwtSecret string) {

	api := app.Group("/api/v1")

	authLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			c.Set("Retry-After", "60")
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "too many attempts, please try again later",
			})
		},
	})

	// public routes
	auth := api.Group("/auth", authLimiter)
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	products := api.Group("/products")
	products.Get("/", productHandler.List)
	products.Get("/:id", productHandler.GetByID)

	// auth required
	authenticated := api.Group("", middleware.Auth(jwtSecret))
	authenticated.Get("/me", userHandler.Me)

	// admin only products
	admin := authenticated.Group("/products", middleware.RequireRole("admin"))
	admin.Post("/", productHandler.Create)
	admin.Put("/:id", productHandler.Update)
	admin.Delete("/:id", productHandler.Delete)

	// carts
	carts := authenticated.Group("/cart")
	carts.Get("/", cartHandler.GetCart)
	carts.Post("/items", cartHandler.AddItem)
	carts.Put("/items/:id", cartHandler.UpdateItem)
	carts.Delete("/items/:id", cartHandler.RemoveItem)

	// orders
	orders := authenticated.Group("/orders")
	orders.Post("/", orderHandler.Checkout)
	orders.Get("/", orderHandler.List)
	orders.Get("/:id", orderHandler.GetByID)
}
