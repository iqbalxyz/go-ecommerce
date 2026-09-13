package middleware

import (
	"go-ecommerce/internal/utils"
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func Auth(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		// 1. Take the authorization header
		authorization := c.Get("Authorization")
		if authorization == "" {
			return utils.Error(c, fiber.StatusUnauthorized, "authorization header is required")
		}

		// 2. Check the format "Bearer <token>"
		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.Error(c, fiber.StatusUnauthorized, "invalid authorization format")
		}
		token := parts[1]

		// 3. Parse token
		claims, err := utils.ParseToken(token, jwtSecret)
		if err != nil {
			log.Printf("Token parsing error: %v", err)
			return utils.Error(c, fiber.StatusUnauthorized, "token is invalid or expired")
		}

		// 4. Set c.Locals("user_id", claims.UserID)
		c.Locals("user_id", claims.UserID)

		// 5. Set c.Locals("role", claims.Role)
		c.Locals("role", claims.Role)

		// 6. return c.Next()
		return c.Next()
	}
}
