package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupCORS mengembalikan instance middleware CORS yang siap digunakan
func SetupCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "http://localhost:6969, http://localhost:1456",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	})
}
