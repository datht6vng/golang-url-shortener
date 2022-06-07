package limiter

import (
	"server_go/config"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func CreateLimiter() fiber.Handler {
	maxRequest := config.Config.Limitter.MaxRequest
	limiterExpire := config.Config.Limitter.LimitterExprire
	intMaxRequest, _ := strconv.Atoi(maxRequest)
	intLimiterExpire, _ := strconv.Atoi(limiterExpire)

	config := limiter.Config{
		Max:        intMaxRequest,
		Expiration: time.Duration(intLimiterExpire) * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).Render("429", nil)
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimiterMiddleware:      limiter.FixedWindow{},
	}
	return limiter.New(config)
}
