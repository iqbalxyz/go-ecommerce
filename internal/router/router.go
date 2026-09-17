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

	// admin only
	admin := authenticated.Group("/products", middleware.RequireRole("admin"))
	admin.Post("/", productHandler.Create)
	admin.Put("/:id", productHandler.Update)
	admin.Delete("/:id", productHandler.Delete)

}
