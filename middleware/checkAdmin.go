package middleware

import (
	"net/http"

	usermodels "github.com/KingSupermarket/model/userModels"
	"github.com/gofiber/fiber/v2"
)

func CheckAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var userAdmin usermodels.Users
		if err := c.BodyParser(&userAdmin); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "faild parser user",
			})
		}

		if userAdmin.Role != "ADMIN" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "You do not have permission to access the data",
			})
		}
		return c.Next()
	}
}
