package controller

import (
	"net/url"
	"server_go/model"

	"github.com/gofiber/fiber/v2"
)

var baseUrl = "http://localhost"

func GetIndexController(ctx *fiber.Ctx) error {
	return ctx.Render("index", fiber.Map{})
}
func PostGenUrlController(ctx *fiber.Ctx) error {
	requestData := model.Url{}
	if err := ctx.BodyParser(&requestData); err != nil {
		return ctx.JSON(fiber.Map{"success": false, "url": nil})
	}
	// check is a valid url
	_, err := url.ParseRequestURI(requestData.Url)
	if err != nil {
		return ctx.JSON(fiber.Map{"success": false, "url": nil})
	}
	// find in database
	// need to find in catch ?
	return nil
}
func GetUrlController(ctx *fiber.Ctx) error {
	requestData := model.Url{}
	if err := ctx.BodyParser(&requestData); err != nil {
		return ctx.JSON(fiber.Map{"success": false, "url": nil})
	}
	result, err := model.FindLongUrl(model.Connection, requestData.Url)
	if err != nil {
		ctx.JSON(err.Error())
	}
	return ctx.JSON(result)
}
func InitDBController(ctx *fiber.Ctx) error {
	err := model.CreateModel(model.Connection)
	if err != nil {
		return ctx.JSON(fiber.Map{"success": false})
	}
	return ctx.JSON(fiber.Map{"success": true})
}
