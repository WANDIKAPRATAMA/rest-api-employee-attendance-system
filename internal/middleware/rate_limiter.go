package middleware

import (
	utils "employee-attendance-system/internal/util"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// SetupRateLimiter mengembalikan instance middleware rate limiter
func SetupRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		// Next: func(c *fiber.Ctx) bool {
		// 	return c.IP() == "127.0.0.1"
		// },
		Max:        50,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			if forwarded := c.Get("x-forwarded-for"); forwarded != "" {
				return forwarded
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(utils.ErrorResponse(
				fiber.StatusTooManyRequests,
				"Route not found",
				[]utils.ErrorDetail{{
					Field:   "Too many requests",
					Message: "Please wait a moment before making another request",
				}},
			))
		},
	})
}
