package controller

import (
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
}
type Url struct {
	User string `json:"user" xml:"user" form:"user"`
	Url  string `json:"url" xml:"url" form:"url"`
}
type Response struct {
	Code    int         `json:"code" xml:"code" form:"code"`
	Message string      `json:"message" xml:"message" form:"message"`
	Data    interface{} `json:"data" xml:"data" form:"data"`
}

func (this *Controller) Success(ctx *fiber.Ctx, successCode int, code int, message string, data interface{}) error {
	response := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	return ctx.Status(successCode).JSON(response)
}
func (this *Controller) Failure(ctx *fiber.Ctx, errorCode int, code int, message string) error {
	response := Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
	return ctx.Status(errorCode).JSON(response)
}
