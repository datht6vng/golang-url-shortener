package middleware

import (
	"trueid-shorten-link/internal/shorten-link/service"

	"github.com/gofiber/fiber/v2"
)

type ValidateAPIKeyMiddleware struct {
	validateAPIKeyService *service.ValidateAPIKeyService
}

func (this *ValidateAPIKeyMiddleware) Init(validateAPIKeyService *service.ValidateAPIKeyService) *ValidateAPIKeyMiddleware {
	this.validateAPIKeyService = validateAPIKeyService
	return this
}

func (this *ValidateAPIKeyMiddleware) ValidateAPIKey(ctx *fiber.Ctx) error {
	apiKey := ctx.Get("X-API-KEY")
	if apiKey == "" {
		return &fiber.Error{
			Code:    401,
			Message: "X-API-KEY header not found!",
		}
	}
	clientID, err := this.validateAPIKeyService.ValidateAPIKey(apiKey)
	if err != nil {
		return err
	}
	ctx.Locals("CLIENT-ID", clientID)
	return ctx.Next()
}
