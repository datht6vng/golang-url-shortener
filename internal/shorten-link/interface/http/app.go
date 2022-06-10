package http

import (
	"github.com/gofiber/fiber/v2"
)

func NewApp() *fiber.App {
	return new(Handler).InitHandler().App
}
