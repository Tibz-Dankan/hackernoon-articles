package health

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

var CheckHealth = func(c *fiber.Ctx) error {
	// Basic checks
	c.Set("Content-Type", "text/plain")
	log.Println("Health checked bro")

	return c.Status(fiber.StatusOK).SendString("OK")
}
