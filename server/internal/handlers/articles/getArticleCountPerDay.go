package articles

import (
	"log"
	"time"

	"github.com/Tibz-Dankan/hackernoon-articles/internal/models"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/pkg"
	"github.com/gofiber/fiber/v2"
)

var GetArticleCountPerDay = func(c *fiber.Ctx) error {
	articles := models.Article{}
	limitParam := c.Query("limit")
	dateCursorParam := c.Query("dateCursor")

	limit, err := pkg.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var parsedDateCursorParam time.Time
	log.Printf("dateCursorParam: %v\n", dateCursorParam)

	if dateCursorParam != "" {
		parsedDateCursorParam, err = time.Parse("2006-01-02", dateCursorParam)
		if err != nil {
			parsedDateCursorParam, err = time.Parse(time.RFC3339, dateCursorParam)
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "Invalid date format! Must be YYYY-MM-DD or ISO 8601 string.")
			}
		}
		log.Printf("parsedDateCursorParam: %v\n", parsedDateCursorParam)
	}

	articleCountPerDay, err := articles.GetArticleCountPerDay(int(limit), parsedDateCursorParam)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var nextCursor string
	if len(articleCountPerDay) > 0 {
		lastDay := articleCountPerDay[len(articleCountPerDay)-1]
		if dateStr, ok := lastDay["date"].(string); ok {
			nextCursor = dateStr
		}
	}

	var totalDays int64
	if err := articles.CountDistinctDays(&totalDays); err != nil {
		log.Printf("Error getting total days count: %v", err)
		totalDays = 0
	}

	count := len(articleCountPerDay)

	pagination := map[string]interface{}{
		"limit":      limit,
		"nextCursor": nextCursor,
		"totalDays":  totalDays,
		"count":      count,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       articleCountPerDay,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
