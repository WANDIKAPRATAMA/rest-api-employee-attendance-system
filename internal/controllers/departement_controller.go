// department_controller.go
package controller

import (
	"employee-attendance-system/internal/entity/dto"
	"employee-attendance-system/internal/middleware"
	"employee-attendance-system/internal/usecase"
	utils "employee-attendance-system/internal/util"
	"math"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type DepartmentController interface {
	CreateDepartment(c *fiber.Ctx) error
	GetDepartment(c *fiber.Ctx) error
	UpdateDepartment(c *fiber.Ctx) error
	DeleteDepartment(c *fiber.Ctx) error
	GetDepartments(c *fiber.Ctx) error // List with pagination
}

type departmentController struct {
	usecase  usecase.DepartmentUseCase
	log      *logrus.Logger
	validate *validator.Validate
}

func NewDepartmentController(usecase usecase.DepartmentUseCase, log *logrus.Logger, validate *validator.Validate) DepartmentController {
	return &departmentController{usecase: usecase, log: log, validate: validate}
}

func (c *departmentController) CreateDepartment(ctx *fiber.Ctx) error {
	var req dto.CreateDepartmentRequest
	allowedFields := utils.GenerateAllowedFields(dto.CreateDepartmentRequest{})
	if err := utils.BindAndValidateBody(ctx, &req, allowedFields, c.validate); err != nil {
		var errors []utils.ErrorDetail
		if validationErr := c.validate.Struct(req); validationErr != nil {
			for _, e := range validationErr.(validator.ValidationErrors) {
				errors = append(errors, utils.ErrorDetail{Field: e.Field(), Message: e.Error()})
			}
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), errors))
	}

	// Check role admin dari middleware (asumsi set di locals)
	if middleware.GetLocalKeys(ctx).Role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(utils.ErrorResponse(fiber.StatusForbidden, "Admin only", nil))
	}

	dept, err := c.usecase.CreateDepartment(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessResponse(fiber.StatusCreated, "Department created", dept, struct{}{}))
}

func (c *departmentController) GetDepartment(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Invalid ID", nil))
	}

	dept, err := c.usecase.GetDepartment(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(fiber.StatusNotFound, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Department retrieved", dept, struct{}{}))
}

func (c *departmentController) UpdateDepartment(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Invalid ID", nil))
	}

	var req dto.UpdateDepartmentRequest
	allowedFields := utils.GenerateAllowedFields(dto.UpdateDepartmentRequest{})
	if err := utils.BindAndValidateBody(ctx, &req, allowedFields, c.validate); err != nil {
		var errors []utils.ErrorDetail
		if validationErr := c.validate.Struct(req); validationErr != nil {
			for _, e := range validationErr.(validator.ValidationErrors) {
				errors = append(errors, utils.ErrorDetail{Field: e.Field(), Message: e.Error()})
			}
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), errors))
	}

	if middleware.GetLocalKeys(ctx).Role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(utils.ErrorResponse(fiber.StatusForbidden, "Admin only", nil))
	}

	dept, err := c.usecase.UpdateDepartment(ctx.Context(), id, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Department updated", dept, struct{}{}))
}

func (c *departmentController) DeleteDepartment(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Invalid ID", nil))
	}

	if middleware.GetLocalKeys(ctx).Role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(utils.ErrorResponse(fiber.StatusForbidden, "Admin only", nil))
	}

	err = c.usecase.DeleteDepartment(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Department deleted", nil, struct{}{}))
}

func (c *departmentController) GetDepartments(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	depts, total, err := c.usecase.GetDepartments(ctx.Context(), page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	pagination := utils.Pagination{
		CurrentPage: page,
		TotalItems:  int(total),
		TotalPages:  int(math.Ceil(float64(total) / float64(limit))),
		HasNextPage: page*limit < int(total),
		NextPage: func() *int {
			if page*limit < int(total) {
				np := page + 1
				return &np
			}
			return nil
		}(),
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Departments retrieved", depts, pagination))
}
