package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type MetricsController struct{}

func (this *MetricsController) Metrics(ctx *fiber.Ctx) error {
	fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())(ctx.Context())
	return nil
}
