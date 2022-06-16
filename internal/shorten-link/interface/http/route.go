package http

func (this *Handler) InitRoute() *Handler {
	this.App.Post("/shorten", this.ValidateAPIKeyMiddleware.ValidateAPIKey, this.GenerateURLController.GenerateURL)
	this.App.Get("/shorten/metrics", this.MetricsController.Metrics)
	this.App.Get("/call/:url", this.ValidateURLMiddleware.ValidateURL, this.GetURLController.GetURL)
	return this
}
