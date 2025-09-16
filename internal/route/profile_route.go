package routes

import (
	controller "employee-attendance-system/internal/controllers"
	"employee-attendance-system/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type UserRouteConfig struct {
	App            *fiber.App
	UserController controller.UserController
	AuthMiddleware *middleware.AuthMiddleware
}

func (r *UserRouteConfig) Setup() {
	api := r.App.Group("/api/v1")
	users := api.Group("/users")
	users.Get("", r.AuthMiddleware.Authenticate, r.UserController.ListUsers)
	profile := api.Group("/profile/")
	profile.Get("", r.AuthMiddleware.Authenticate, r.UserController.GetProfile)    // PUT /api/v1/profile
	profile.Put("", r.AuthMiddleware.Authenticate, r.UserController.UpdateProfile) // PUT /api/v1/profile
}
