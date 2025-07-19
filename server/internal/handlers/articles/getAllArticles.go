package articles

import (
	"github.com/Tibz-Dankan/hackernoon-articles/internal/models"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/pkg"
	"github.com/gofiber/fiber/v2"
)

// To significantly improve the article
var GetAllArticles = func(c *fiber.Ctx) error {
	articles := models.Article{}
	limitParam := c.Query("limit")
	cursorParam := c.Query("cursor")

	limit, err := pkg.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if cursorParam == "" {
		cursorParam = ""
	}

	allArticles, err := articles.FindAllByPostedAt(limit, cursorParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var prevCursor string
	if len(allArticles) > 0 {
		prevCursor = allArticles[len(allArticles)-1].ID
	}

	pagination := map[string]interface{}{
		"limit":      limit,
		"prevCursor": prevCursor,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       allArticles,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
