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
	clientID := ctx.Locals("CLIENT-ID").(string)
	//clientLimit := ctx.Locals("CLIENT-LIMIT").(int64)
	clientLimit := int64(10000)
	// take limit
	err := this.limitGenerateService.LimitGenerate(clientID, clientLimit)
	if err != nil {
		return err
	}
	return ctx.Next()
}
