package status

import (
	"github.com/gofiber/fiber/v2"
)

var GetAppStatus = func(c *fiber.Ctx) error {
	response := fiber.Map{
		"status":  "success",
		"message": "Active",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
