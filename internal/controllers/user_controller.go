package controller

import (
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

type UserController interface {
	ListUsers(c *fiber.Ctx) error
	UpdateProfile(ctx *fiber.Ctx) error
	GetProfile(ctx *fiber.Ctx) error
}

type userController struct {
	usecase  usecase.UserUseCase
	log      *logrus.Logger
	validate *validator.Validate
}

func NewUserController(usecase usecase.UserUseCase, log *logrus.Logger, validate *validator.Validate) UserController {
	return &userController{usecase: usecase, log: log, validate: validate}
}

func (c *userController) ListUsers(ctx *fiber.Ctx) error {
	var req dto.ListUsersRequest
	req.Page = ctx.QueryInt("page", 1)
	req.Limit = ctx.QueryInt("limit", 10)
	req.Email = ctx.Query("email")
	req.Status = ctx.Query("status")
	if depID := ctx.Query("department_id"); depID != "" {
		id, err := uuid.Parse(depID)
		if err == nil {
			req.DepartmentID = &id
		}
	}
	if start := ctx.Query("created_at_start"); start != "" {
		t, err := time.Parse("2006-01-02", start)
		if err == nil {
			req.CreatedAtStart = &t
		}
	}
	if end := ctx.Query("created_at_end"); end != "" {
		t, err := time.Parse("2006-01-02", end)
		if err == nil {
			req.CreatedAtEnd = &t
		}
	}

	if err := c.validate.Struct(req); err != nil {
		var errors []utils.ErrorDetail
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, utils.ErrorDetail{Field: e.Field(), Message: e.Error()})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Validation failed", errors))
	}

	// Check role
	role := middleware.GetLocalKeys(ctx).Role
	if role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(utils.ErrorResponse(fiber.StatusForbidden, "Admin only", nil))
	}

	users, total, err := c.usecase.ListUsers(ctx.Context(), req)
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

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Users retrieved", users, pagination))
}
func (c *userController) UpdateProfile(ctx *fiber.Ctx) error {
	var req dto.UpdateProfileRequest
	allowedFields := utils.GenerateAllowedFields(dto.UpdateProfileRequest{})
	if err := utils.BindAndValidateBody(ctx, &req, allowedFields, c.validate); err != nil {
		var errors []utils.ErrorDetail
		if validationErr := c.validate.Struct(req); validationErr != nil {
			for _, e := range validationErr.(validator.ValidationErrors) {
				var msg string
				switch e.Tag() {
				case "phone":
					msg = "Phone number must be 8-15 digits (optionally with +)"
				case "url":
					msg = "Invalid URL format"
				default:
					msg = e.Error()
				}

				errors = append(errors, utils.ErrorDetail{
					Field:   e.Field(),
					Message: msg,
				})
			}
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), errors))
	}

	userID := middleware.GetLocalKeys(ctx).UserID

	updatedProfile, err := c.usecase.UpdateProfile(ctx.Context(), userID, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Profile updated successfully", fiber.Map{
		"id":            updatedProfile.ID,
		"full_name":     updatedProfile.FullName,
		"phone":         updatedProfile.Phone,
		"avatar_url":    updatedProfile.AvatarURL,
		"address":       updatedProfile.Address,
		"department_id": updatedProfile.DepartmentID,
	}, nil))
}

func (c *userController) GetProfile(ctx *fiber.Ctx) error {
	userID := middleware.GetLocalKeys(ctx).UserID

	profile, err := c.usecase.GetProfile(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(
		fiber.StatusCreated,
		"Profile Retrieved Successfully",
		profile,
		struct{}{},
	))
}
