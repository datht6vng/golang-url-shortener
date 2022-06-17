package controller

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type ErrorController struct {
	Controller
}

func (this *ErrorController) ErrorController(ctx *fiber.Ctx, err error) error {
	// Default 500 statuscode
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		// Override status code if fiber.Error type
		code = e.Code
	}
	// For API path
	if len(ctx.Path()) > 4 && ctx.Path()[:4] == "/api" && ctx.Path()[:4] == "/api" {
		if code >= 500 {
			log.Println(err.Error())
		}
		this.Failure(ctx, code, code, err.Error())
		return nil
	}
	// Set Content-Type: text/plain; charset=utf-8
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	// Return statuscode with error message
	// Log internal server error
	if code >= 500 {
		log.Println(err.Error())
	}
	return ctx.Status(code).Render(fmt.Sprint(code), fiber.Map{"error": err.Error()})
}
