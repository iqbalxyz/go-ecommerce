package router

import (
	"go-ecommerce/internal/handler"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler) {
	// create group /api/v1
	api := app.Group("/api/v1")

	// auth group
	auth := api.Group("/auth")

	// user auth routes
	auth.Post("/register", userHandler.Register)
}
