package routes

import (
	controller "employee-attendance-system/internal/controllers"
	"employee-attendance-system/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	AuthController controller.AuthController
	AuthMiddleware *middleware.AuthMiddleware
}

func (r *RouteConfig) Setup() {
	api := r.App.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/signup", r.AuthController.Signup)
	auth.Post("/signin", r.AuthController.Signin)
	auth.Post("/change-password", r.AuthMiddleware.Authenticate, r.AuthController.ChangePassword)
	auth.Post("/refresh-token", r.AuthController.RefreshToken)
	auth.Post("/change-role", r.AuthMiddleware.Authenticate, r.AuthController.ChangeRole)
	auth.Post("/signout", r.AuthMiddleware.Authenticate, r.AuthController.Signout)

}
