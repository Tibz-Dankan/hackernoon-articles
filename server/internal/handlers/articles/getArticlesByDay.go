package articles

import (
	"log"
	"time"

	"github.com/Tibz-Dankan/hackernoon-articles/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetArticlesByDay = func(c *fiber.Ctx) error {
	articles := models.Article{}
	postedAtParam := c.Params("postedAt")
	var parsedPostedAtParam time.Time
	var err error

	log.Printf("postedAtParam: %v\n", postedAtParam)
	if postedAtParam != "" {
		parsedPostedAtParam, err = time.Parse(time.RFC3339, postedAtParam)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid offset format! Must be an ISO 8601 string.")
		}
		log.Printf("parsedDateCursorParam: %v\n", parsedPostedAtParam)
	}

	allArticles, err := articles.FindByPostedAt(parsedPostedAtParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":     "success",
		"data":       allArticles,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}