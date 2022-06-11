package controller

import (
	"net/http"
	"trueid-shorten-link/internal/shorten-link/service"

	"github.com/gofiber/fiber/v2"
)

type GenerateURLController struct {
	Controller
	generateURLService *service.GenerateURLService
}
type RequestData struct {
	URL string `json:"url" xml:"url" form:"url"`
}

func (this *GenerateURLController) Init(generateURLService *service.GenerateURLService) *GenerateURLController {
	this.generateURLService = generateURLService
	return this
}
func (this *GenerateURLController) GenerateURL(ctx *fiber.Ctx) error {
	requestData := new(RequestData)
	if err := ctx.BodyParser(requestData); err != nil {
		return this.Failure(ctx, http.StatusBadRequest, http.StatusBadRequest, err.Error())
	}
	shortURL, err := this.generateURLService.GenerateURL(requestData.URL, ctx.Locals("CLIENT-ID").(string))
	if err != nil {
		return err
	}
	return this.Success(ctx, http.StatusOK, http.StatusOK, "Success!", &fiber.Map{"url": ctx.BaseURL() + "/call/" + shortURL})
}
