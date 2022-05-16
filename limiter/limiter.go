package limiter

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func CreateLimiter() fiber.Handler {
	maxRequest := os.Getenv("MAX_REQUEST")
	limitExpire := os.Getenv("LIMIT_EXPIRE")

	if maxRequest == "" {
		maxRequest = "5"
	}
	if limitExpire == "" {
		limitExpire = "1"
	}
	intMaxRequest, _ := strconv.Atoi(maxRequest)
	intLimitExpire, _ := strconv.Atoi(limitExpire)

	config := limiter.Config{
		Max:        intMaxRequest,
		Expiration: time.Duration(intLimitExpire) * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimiterMiddleware:      limiter.FixedWindow{},
	}
	return limiter.New(config)
}
