package router

import (
	"go-ecommerce/internal/handler"
	"go-ecommerce/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler, jwtSecret string) {
	// create group /api/v1
	api := app.Group("/api/v1")

	// auth group
	auth := api.Group("/auth")

	// user auth routes
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	// users route group
	users := api.Group("/users", middleware.Auth(jwtSecret))
	users.Get("/me", userHandler.Me)
}
