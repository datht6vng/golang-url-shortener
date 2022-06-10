package controller

import (
	"net/http"
	"trueid-shorten-link/internal/shorten-link/service"

	"github.com/gofiber/fiber/v2"
)

type GenerateUrlController struct {
	Controller
	generateUrlService *service.GenerateUrlService
}

func (this *GenerateUrlController) Init(generateUrlService *service.GenerateUrlService) *GenerateUrlController {
	this.generateUrlService = generateUrlService
	return this
}
func (this *GenerateUrlController) GenerateUrl(ctx *fiber.Ctx) error {
	inputData := new(service.UrlData)
	if err := ctx.BodyParser(inputData); err != nil {
		return this.Failure(ctx, http.StatusBadRequest, http.StatusBadRequest, err.Error())
	}
	shortUrl, err := this.generateUrlService.GenerateUrl(inputData)
	if err != nil {
		return err
	}
	return this.Success(ctx, http.StatusOK, http.StatusOK, "Success!", &fiber.Map{"url": ctx.BaseURL() + "/call/" + shortUrl})
}
