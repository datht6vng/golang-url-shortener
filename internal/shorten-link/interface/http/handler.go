package http

import (
	"time"
	"trueid-shorten-link/config"
	"trueid-shorten-link/internal/shorten-link/interface/http/controller"
	"trueid-shorten-link/internal/shorten-link/interface/job"
	"trueid-shorten-link/internal/shorten-link/repository"
	"trueid-shorten-link/internal/shorten-link/service"
	"trueid-shorten-link/migration"
	"trueid-shorten-link/package/database"
	"trueid-shorten-link/package/metrics"
	"trueid-shorten-link/package/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
)

type Handler struct {
	App                   *fiber.App
	ErrorController       *controller.ErrorController
	GenerateUrlController *controller.GenerateUrlController
	GetUrlController      *controller.GetUrlController
	MetricsController     *controller.MetricsController
}

func (this *Handler) InitHandler() *Handler {
	// Connect database
	db := database.Connect()
	migration.CreateTable(db)
	redis := redis.Connect()

	// Metrics
	metrics := new(metrics.Metrics).Init()

	// Repositories
	urlRepository := new(repository.UrlRepository).Init(db)

	// Job
	job := new(job.Job).Init(urlRepository, redis)
	job.CreateCronJob(
		"@every 12h",
		job.DeleteExpireUrl,
		job.ResetMaxID,
	)
	job.DeleteExpireUrl()
	job.ResetMaxID()

	// Controllers
	this.ErrorController = new(controller.ErrorController)
	this.GenerateUrlController = new(controller.GenerateUrlController).Init(
		new(service.GenerateUrlService).Init(urlRepository, redis, metrics),
	)
	this.GetUrlController = new(controller.GetUrlController).Init(
		new(service.GetUrlService).Init(urlRepository, redis, metrics),
	)

	viewEngine := html.New(config.Config.View.Path, ".html")

	this.App = fiber.New(fiber.Config{
		Views:        viewEngine,
		ErrorHandler: this.ErrorController.ErrorController,
	})

	this.App.Use(fiberRecover.New())

	this.App.Use(limiter.New(limiter.Config{
		Max:        config.Config.Limitter.MaxRequest,
		Expiration: time.Duration(config.Config.Limitter.LimitterExprire) * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).Render("429", nil)
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimiterMiddleware:      limiter.FixedWindow{},
	}))

	this.InitRoute()
	return this
}
