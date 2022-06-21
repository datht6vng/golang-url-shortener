package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// NewLogMiddleWare creates a new middleware handler
func NewLogMiddleWare(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Ignore path
		url := c.OriginalURL()
		if url == "/health" || url == "/api/metrics" {
			return c.Next()
		}
		var body string
		if url == "/api/update_user" || url == "/api/agents/create" || url == "/api/login" || url == "api/authorize" {
			body = ""
		} else {
			body = string(c.Body())
		}

		var start, stop time.Time

		// Set latency start time
		start = time.Now()

		// Handle request, store err for logging
		err := c.Next()
		stop = time.Now()

		logger.Infof("%s %s With body:%s response:%s executeTime:%dms ", c.Method(), c.OriginalURL(), body, string(c.Response().Body()), stop.Sub(start)/time.Millisecond)

		return err
	}
}
