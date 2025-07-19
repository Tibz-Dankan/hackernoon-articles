package pkg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var DefaultErrorHandler = func(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var status = "error"
	var message = "Ooops, something went wrong, try again later!"
	var detail = err.Error()

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		if strings.HasPrefix(strconv.Itoa(code), "4") {
			message = err.Error()
			status = "fail"
		}
	}

	errDetailMsg := struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Detail  string `json:"detail"`
	}{
		Code:    code,
		Status:  status,
		Message: message,
		Detail:  detail,
	}
	fmt.Println("ERROR💥:", errDetailMsg)

	errMsg := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  status,
		Message: message,
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(code).JSON(errMsg)
}
