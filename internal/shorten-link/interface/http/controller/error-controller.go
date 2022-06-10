package controller

import (
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
	if ctx.Path()[:4] == "/api" {
		if code >= 500 {
			log.Println(err.Error())
		}
		this.Failure(ctx, code, code, err.Error())
		return nil
	}
	// Set Content-Type: text/plain; charset=utf-8
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	// Return statuscode with error message
	if code == 404 {
		return ctx.Status(404).Render("404", nil)
	}
	// Log internal server error
	if code >= 500 {
		log.Println(err.Error())
	}
	return ctx.Status(500).Render("500", fiber.Map{"err": err.Error()})
}
