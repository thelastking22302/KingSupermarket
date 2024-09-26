package middleware

import (
	"net/http"
	"strings"

	"github.com/KingSupermarket/pkg/security"
	"github.com/gofiber/fiber/v2"
)

func JwtMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header is missing"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := security.ValidateToken(tokenString)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		// Store userId in context for use in handlers
		c.Locals("userId", claims.User_Id)
		return c.Next()
	}
}
