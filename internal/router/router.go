package router

import (
	"go-ecommerce/internal/handler"
	"go-ecommerce/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(
	app *fiber.App,
	userHandler *handler.UserHandler,
	productHandler *handler.ProductHandler,
	cartHandler *handler.CartHandler,
	orderHandler *handler.OrderHandler,
	jwtSecret string) {

	api := app.Group("/api/v1")

	// public routes
	auth := api.Group("/auth")
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
