package http

func (this *Handler) InitRoute() *Handler {
	this.App.Post("/api/generate-url",
		this.ValidateAPIKeyMiddleware.ValidateAPIKey,
		this.LimitGenerateMiddleware.LimitGenerate,
		this.GenerateURLController.GenerateURL,
	)
	this.App.Get("/api/metrics", this.MetricsController.Metrics)
	this.App.Post("/api/client", this.ClientController.UpdateClient)
	this.App.Get("/:url", this.ValidateURLMiddleware.ValidateURL, this.GetURLController.GetURL)
	return this
}
