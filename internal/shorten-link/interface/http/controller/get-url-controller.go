package controller

import (
	"trueid-shorten-link/internal/shorten-link/service"

	"github.com/gofiber/fiber/v2"
)

type GetUrlController struct {
	Controller
	getUrlService *service.GetUrlService
}

func (this *GetUrlController) Init(getUrlService *service.GetUrlService) *GetUrlController {
	this.getUrlService = getUrlService
	return this
}
func (this *GetUrlController) GetUrl(ctx *fiber.Ctx) error {
	url := ctx.Params("url")
	longUrl, err := this.getUrlService.GetUrl(url)
	if err != nil {
		return err
	}
	return ctx.Redirect(longUrl)
}
