package middleware

import "github.com/gofiber/fiber/v2"

func ValidateUrl(ctx *fiber.Ctx) error {
	return ctx.Next()
}
