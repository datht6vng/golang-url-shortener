package controller

import (
	"trueid-shorten-link/internal/shorten-link/service"

	"github.com/gofiber/fiber/v2"
)

type GetURLController struct {
	Controller
	getURLService *service.GetURLService
}

func (this *GetURLController) Init(getURLService *service.GetURLService) *GetURLController {
	this.getURLService = getURLService
	return this
}
func (this *GetURLController) GetURL(ctx *fiber.Ctx) error {
	url := ctx.Params("url")
	longURL, err := this.getURLService.GetURL(url)
	if err != nil {
		return err
	}
	return ctx.Redirect(longURL)
}
