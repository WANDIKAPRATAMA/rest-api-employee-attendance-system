package routes

import (
	controller "employee-attendance-system/internal/controllers"
	"employee-attendance-system/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type DepartmentRouteConfig struct {
	App                  *fiber.App
	DepartmentController controller.DepartmentController
	AuthMiddleware       *middleware.AuthMiddleware
}

func (r *DepartmentRouteConfig) Setup() {
	api := r.App.Group("/api/v1")
	dept := api.Group("/departments")
	dept.Post("", r.AuthMiddleware.Authenticate, r.DepartmentController.CreateDepartment)
	dept.Get("/:id", r.AuthMiddleware.Authenticate, r.DepartmentController.GetDepartment)
	dept.Put("/:id", r.AuthMiddleware.Authenticate, r.DepartmentController.UpdateDepartment)
	dept.Delete("/:id", r.AuthMiddleware.Authenticate, r.DepartmentController.DeleteDepartment)
	dept.Get("", r.AuthMiddleware.Authenticate, r.DepartmentController.GetDepartments) // List
	dept.Post("/assignment", r.AuthMiddleware.Authenticate, r.DepartmentController.AssignmentDepartement)

}
