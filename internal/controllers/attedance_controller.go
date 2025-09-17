// attendance_controller.go
package controller

import (
	"employee-attendance-system/internal/entity/domain"
	"employee-attendance-system/internal/entity/dto"
	"employee-attendance-system/internal/middleware"
	"employee-attendance-system/internal/usecase"
	utils "employee-attendance-system/internal/util"
	"math"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AttendanceController interface {
	ClockIn(c *fiber.Ctx) error
	ClockOut(c *fiber.Ctx) error
	GetAttendanceLogs(c *fiber.Ctx) error
	GetAdminDashboard(ctx *fiber.Ctx) error
	GetAttendanceHistory(ctx *fiber.Ctx) error
	CheckCurrentStatus(ctx *fiber.Ctx) error
}

type attendanceController struct {
	usecase  usecase.AttendanceUseCase
	log      *logrus.Logger
	validate *validator.Validate
}

func NewAttendanceController(usecase usecase.AttendanceUseCase, log *logrus.Logger, validate *validator.Validate) AttendanceController {
	return &attendanceController{usecase: usecase, log: log, validate: validate}
}
func (c *attendanceController) GetAttendanceHistory(ctx *fiber.Ctx) error {
	var req dto.GetAttendanceHistoryRequest
	req.Page = ctx.QueryInt("page", 1)
	req.Limit = ctx.QueryInt("limit", 10)
	userID, err := uuid.Parse(ctx.Query("user_id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Invalid user_id", nil))
	}
	req.UserID = userID

	if err := c.validate.Struct(req); err != nil {
		var errors []utils.ErrorDetail
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, utils.ErrorDetail{Field: e.Field(), Message: e.Error()})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Validation failed", errors))
	}

	// Cek apakah userID milik pengguna saat ini atau admin
	currentUserID := middleware.GetLocalKeys(ctx).UserID
	c.log.Info("🚀 ~ currentUserID:", currentUserID)
	if req.UserID != currentUserID && middleware.GetLocalKeys(ctx).Role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(utils.ErrorResponse(fiber.StatusForbidden, "Access denied", nil))
	}

	histories, total, err := c.usecase.GetAttendanceHistory(ctx.Context(), req)
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

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Attendance history retrieved", histories, pagination))
}
func (c *attendanceController) GetAdminDashboard(ctx *fiber.Ctx) error {
	var req dto.AdminDashboardRequest
	if start := ctx.Query("start_date"); start != "" {
		t, err := time.Parse("2006-01-02", start)
		if err == nil {
			req.StartDate = &t
		}
	}
	if end := ctx.Query("end_date"); end != "" {
		t, err := time.Parse("2006-01-02", end)
		if err == nil {
			req.EndDate = &t
		}
	}

	if err := c.validate.Struct(req); err != nil {
		var errors []utils.ErrorDetail
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, utils.ErrorDetail{Field: e.Field(), Message: e.Error()})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Validation failed", errors))
	}

	if middleware.GetLocalKeys(ctx).Role != "admin" {
		c.log.Printf("role: %v", middleware.GetLocalKeys(ctx).Role)
		return ctx.Status(fiber.StatusForbidden).JSON(utils.ErrorResponse(fiber.StatusForbidden, "Admin only", nil))
	}

	dashboard, err := c.usecase.GetAdminDashboard(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Dashboard data retrieved", dashboard, struct{}{}))
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

// func (c *attendanceController) CheckCurrentStatus(ctx *fiber.Ctx) error {
// 	var req dto.CheckCurrentStatusRequest
// 	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
// 		userID, err := uuid.Parse(userIDStr)
// 		if err != nil {
// 			return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Invalid user_id", nil))
// 		}
// 		req.UserID = &userID
// 	}

// 	if err := c.validate.Struct(req); err != nil {
// 		var errors []utils.ErrorDetail
// 		for _, e := range err.(validator.ValidationErrors) {
// 			errors = append(errors, utils.ErrorDetail{Field: e.Field(), Message: e.Error()})
// 		}
// 		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Validation failed", errors))
// 	}

// 	currentUserID := middleware.GetLocalKeys(ctx).UserID
// 	role := middleware.GetLocalKeys(ctx).Role

// 	// Tentukan userID yang akan dicek
// 	var targetUserID uuid.UUID

// 	if req.UserID != nil {
// 		if req.UserID != nil && *req.UserID != currentUserID && role != string(domain.Admin) {
// 			return ctx.Status(fiber.StatusForbidden).JSON(utils.ErrorResponse(fiber.StatusForbidden, "Admin only for other users", nil))
// 		}
// 		targetUserID = *req.UserID
// 	} else {
// 		targetUserID = currentUserID
// 	}

// 	status, err := c.usecase.CheckCurrentStatus(ctx.Context(), targetUserID)
// 	if err != nil {
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Current status retrieved", status, struct{}{}))
// }

func (c *attendanceController) CheckCurrentStatus(ctx *fiber.Ctx) error {
	var req dto.CheckCurrentStatusRequest

	// Ambil query param user_id
	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).
				JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Invalid user_id", nil))
		}
		req.UserID = &userID
	}

	// Validasi request
	if err := c.validate.Struct(req); err != nil {
		var errors []utils.ErrorDetail
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, utils.ErrorDetail{
				Field:   e.Field(),
				Message: e.Error(),
			})
		}
		return ctx.Status(fiber.StatusBadRequest).
			JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Validation failed", errors))
	}

	// Ambil user dari middleware
	localKeys := middleware.GetLocalKeys(ctx)
	currentUserID := localKeys.UserID
	role := localKeys.Role

	// Tentukan target user
	var targetUserID uuid.UUID
	if req.UserID != nil {
		if *req.UserID != currentUserID && role != string(domain.Admin) {
			return ctx.Status(fiber.StatusForbidden).
				JSON(utils.ErrorResponse(fiber.StatusForbidden, "Admin only for other users", nil))
		}
		targetUserID = *req.UserID
	} else {
		targetUserID = currentUserID
	}

	// Panggil usecase
	status, err := c.usecase.CheckCurrentStatus(ctx.Context(), targetUserID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.SuccessResponse(fiber.StatusOK, "Current status retrieved", status, struct{}{}))
}
