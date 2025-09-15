// attendance_controller.go
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

type AttendanceController interface {
	ClockIn(c *fiber.Ctx) error
	ClockOut(c *fiber.Ctx) error
	GetAttendanceLogs(c *fiber.Ctx) error
}

type attendanceController struct {
	usecase  usecase.AttendanceUseCase
	log      *logrus.Logger
	validate *validator.Validate
}

func NewAttendanceController(usecase usecase.AttendanceUseCase, log *logrus.Logger, validate *validator.Validate) AttendanceController {
	return &attendanceController{usecase: usecase, log: log, validate: validate}
}

func (c *attendanceController) ClockIn(ctx *fiber.Ctx) error {
	userID := middleware.GetLocalKeys(ctx).UserID

	attendance, err := c.usecase.ClockIn(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Clocked in", attendance, struct{}{}))
}

func (c *attendanceController) ClockOut(ctx *fiber.Ctx) error {
	userID := middleware.GetLocalKeys(ctx).UserID

	attendance, err := c.usecase.ClockOut(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Clocked out", attendance, struct{}{}))
}

func (c *attendanceController) GetAttendanceLogs(ctx *fiber.Ctx) error {
	var req dto.GetAttendanceLogsRequest
	req.Page = ctx.QueryInt("page", 1)
	req.Limit = ctx.QueryInt("limit", 10)
	req.Date = ctx.Query("date")
	departmentIDStr := ctx.Query("department_id")
	if departmentIDStr != "" {
		if parsedID, err := uuid.Parse(departmentIDStr); err == nil {
			req.DepartmentID = &parsedID
		}
	}
	if err := c.validate.Struct(req); err != nil {
		var errors []utils.ErrorDetail
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, utils.ErrorDetail{Field: e.Field(), Message: e.Error()})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Validation failed", errors))
	}

	// Check admin for full access; else limit to own department (asumsi)
	role := middleware.GetLocalKeys(ctx).Role
	userID := middleware.GetLocalKeys(ctx).UserID

	logs, total, err := c.usecase.GetAttendanceLogs(ctx.Context(), userID, role, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	pagination := utils.Pagination{
		CurrentPage: req.Page,
		TotalItems:  int(total),
		TotalPages:  int(math.Ceil(float64(total) / float64(req.Limit))),
		HasNextPage: req.Page*req.Limit < int(total),
		NextPage: func() *int {
			if req.Page*req.Limit < int(total) {
				np := req.Page + 1
				return &np
			}
			return nil
		}(),
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Attendance logs retrieved", logs, pagination))
}
