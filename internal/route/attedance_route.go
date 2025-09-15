package routes

import (
	controller "employee-attendance-system/internal/controllers"
	"employee-attendance-system/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type AttendanceRouteConfig struct {
	App                  *fiber.App
	AttendanceController controller.AttendanceController
	AuthMiddleware       *middleware.AuthMiddleware
}

func (r *AttendanceRouteConfig) Setup() {
	api := r.App.Group("/api/v1")
	att := api.Group("/attendance")
	att.Post("/clock-in", r.AuthMiddleware.Authenticate, r.AttendanceController.ClockIn)
	att.Put("/clock-out", r.AuthMiddleware.Authenticate, r.AttendanceController.ClockOut)
	att.Get("/logs", r.AuthMiddleware.Authenticate, r.AttendanceController.GetAttendanceLogs)
}
