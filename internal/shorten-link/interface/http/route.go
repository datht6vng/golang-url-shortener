package http

import (
	"trueid-shorten-link/internal/shorten-link/interface/http/middleware"
)

func (this *Handler) InitRoute() *Handler {
	this.App.Post("/api/generate-url", this.GenerateUrlController.GenerateUrl)
	this.App.Get("/api/metrics", this.MetricsController.Metrics)
	this.App.Get("/call/:url", middleware.ValidateUrl, this.GetUrlController.GetUrl)
	return this
}
