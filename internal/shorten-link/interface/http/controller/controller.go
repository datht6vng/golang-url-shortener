package controller

import (
	"trueid-shorten-link/internal/shorten-link/entity"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
}

func (this *Controller) Success(ctx *fiber.Ctx, successCode int, code int, message string, data interface{}) error {
	response := entity.Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	return ctx.Status(successCode).JSON(response)
}
func (this *Controller) Failure(ctx *fiber.Ctx, errorCode int, code int, message string) error {
	response := entity.Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
	return ctx.Status(errorCode).JSON(response)
}
