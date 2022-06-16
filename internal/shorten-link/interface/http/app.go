package http

import (
	"time"
	"trueid-shorten-link/internal/shorten-link/interface/http/controller"
	"trueid-shorten-link/internal/shorten-link/interface/http/middleware"
	"trueid-shorten-link/internal/shorten-link/interface/job"
	"trueid-shorten-link/internal/shorten-link/repository"
	"trueid-shorten-link/internal/shorten-link/service"
	"trueid-shorten-link/package/config"
	"trueid-shorten-link/package/database"
	"trueid-shorten-link/package/metrics"
	"trueid-shorten-link/package/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
)

func NewApp() *fiber.App {
	return new(Handler).InitHandler().App
}

type Handler struct {
	App                      *fiber.App
	ErrorController          *controller.ErrorController
	GenerateURLController    *controller.GenerateURLController
	GetURLController         *controller.GetURLController
	MetricsController        *controller.MetricsController
	ValidateURLMiddleware    *middleware.ValidateURLMiddleware
	ValidateAPIKeyMiddleware *middleware.ValidateAPIKeyMiddleware
	LimitGenerateMiddleware  *middleware.LimitGenerateMiddleware
}

func (this *Handler) InitHandler() *Handler {
	// Connect database
	db := database.Connect()
	redis := redis.Connect()

	// Metrics
	metrics := new(metrics.Metrics).Init()

	// Repositories
	urlRepository := new(repository.URLRepository).Init(db)
	clientRepository := new(repository.ClientRepository).Init(db)
	generateCounterRepository := new(repository.GenerateCounterRepository).Init(db)
	// Job
	job := new(job.Job).Init(urlRepository, generateCounterRepository, redis)
	job.CreateCronJob(
		"@every 12h",
		job.DeleteExpireURL,
		job.ResetMaxID,
	)
	job.CreateCronJob(
		"@every 30m",
		job.UpdateLinkCounter,
	)
	job.DeleteExpireURL()
	job.ResetMaxID()
	job.UpdateLinkCounter()
	// Controllers
	this.ErrorController = new(controller.ErrorController)
	this.GenerateURLController = new(controller.GenerateURLController).Init(
		new(service.GenerateURLService).Init(urlRepository, redis, metrics),
	)
	this.GetURLController = new(controller.GetURLController).Init(
		new(service.GetURLService).Init(urlRepository, redis, metrics),
	)
	// Middlewares
	this.ValidateURLMiddleware = new(middleware.ValidateURLMiddleware).Init()
	this.ValidateAPIKeyMiddleware = new(middleware.ValidateAPIKeyMiddleware).Init(
		new(service.ValidateAPIKeyService).Init(clientRepository, redis),
	)
	this.LimitGenerateMiddleware = new(middleware.LimitGenerateMiddleware).Init(
		new(service.LimitGenerateService).Init(urlRepository, redis),
	)
	this.App = fiber.New(fiber.Config{
		Views:        html.New(config.Config.View.Path, ".html"),
		ErrorHandler: this.ErrorController.ErrorController,
	})
	this.App.Use(recover.New())
	this.App.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))
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
