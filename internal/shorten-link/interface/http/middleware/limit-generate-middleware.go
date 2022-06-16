package middleware

import (
	"trueid-shorten-link/internal/shorten-link/service"

	"github.com/gofiber/fiber/v2"
)

type LimitGenerateMiddleware struct {
	limitGenerateService *service.LimitGenerateService
}

func (this *LimitGenerateMiddleware) Init(limitGenerateService *service.LimitGenerateService) *LimitGenerateMiddleware {
	this.limitGenerateService = limitGenerateService
	return this
}
func (this *LimitGenerateMiddleware) LimitGenerate(ctx *fiber.Ctx) error {
	err := this.limitGenerateService.LimitGenerate(ctx.Locals("CLIENT-ID").(string), ctx.Locals("MAX-LINK").(int64))
	if err != nil {
		return err
	}
	return ctx.Next()
}
