package articles

import (
	"log"
	"strconv"
	"time"

	"github.com/Tibz-Dankan/hackernoon-articles/internal/models"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/pkg"
	"github.com/gofiber/fiber/v2"
)

var GetAllArticles = func(c *fiber.Ctx) error {
	articles := models.Article{}
	limitParam := c.Query("limit")
	articleIDCursorParam := c.Query("articleIDCursor")
	dateCursorParam := c.Query("dateCursor")
	offsetParam := c.Query("offset")

	limit, err := pkg.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if articleIDCursorParam == "" {
		articleIDCursorParam = ""
	}

	var parsedDateCursorParam time.Time
	var offset int

	log.Printf("dateCursorParam: %v\n", dateCursorParam)
	if dateCursorParam != "" {
		parsedDateCursorParam, err := time.Parse(time.RFC3339, dateCursorParam)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid offset format! Must be an ISO 8601 string.")
		}
		log.Printf("parsedDateCursorParam: %v\n", parsedDateCursorParam)
	}

	if offsetParam != "" {
		offset, err = strconv.Atoi(offsetParam)
		if err != nil {
			log.Println("Error converting offsetParam to an integer:", err)
		}
		log.Println(offset)
	}

	allArticles, count, err := articles.FindAllByPostedAt(int(limit), articleIDCursorParam, parsedDateCursorParam, offset)
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
		"count":      count,
		"offset":     offset,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       allArticles,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
