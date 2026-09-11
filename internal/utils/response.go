package utils

import "github.com/gofiber/fiber/v3"

func Success(c fiber.Ctx, status int, data any) error {
	response := fiber.Map{
		"success": true,
		"data":    data,
	}

	return c.Status(status).JSON(response)
}

func Error(c fiber.Ctx, status int, message string) error {
	response := fiber.Map{
		"success": false,
		"error":   message,
	}

	return c.Status(status).JSON(response)
}

func ValidationError(c fiber.Ctx, errors map[string]string) error {
	response := fiber.Map{
		"success": false,
		"error":   "validation failed",
		"errors":  errors,
	}

	return c.Status(400).JSON(response)
}
