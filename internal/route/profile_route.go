package routes

import (
	controller "employee-attendance-system/internal/controllers"
	"employee-attendance-system/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type ProfileRouteConfig struct {
	App            *fiber.App
	AuthController controller.AuthController
	AuthMiddleware *middleware.AuthMiddleware
}

func (r *ProfileRouteConfig) Setup() {
	api := r.App.Group("/api/v1")

	profile := api.Group("/profile/")
	profile.Get("", r.AuthMiddleware.Authenticate, r.AuthController.GetProfile)    // PUT /api/v1/profile
	profile.Put("", r.AuthMiddleware.Authenticate, r.AuthController.UpdateProfile) // PUT /api/v1/profile
}
