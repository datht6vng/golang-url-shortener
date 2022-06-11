package middleware

import (
	"net/http"
	"trueid-shorten-link/package/encryption"

	"github.com/gofiber/fiber/v2"
)

type ValidateURLMiddleware struct{}

func (this *ValidateURLMiddleware) Init() *ValidateURLMiddleware { return this }
func (this *ValidateURLMiddleware) ValidateURL(ctx *fiber.Ctx) error {
	url := ctx.Params("url")
	if len(url) < 5 {
		return &fiber.Error{
			Code:    http.StatusNotFound,
			Message: "URL not found!",
		}
	}
	urlPart := url[:len(url)-4]
	userSignature := url[len(url)-4:]
	systemSignature := encryption.Signature(urlPart)[:4]
	if userSignature != systemSignature {
		return &fiber.Error{
			Code:    http.StatusNotFound,
			Message: "URL not found!",
		}
	}
	return ctx.Next()
}
