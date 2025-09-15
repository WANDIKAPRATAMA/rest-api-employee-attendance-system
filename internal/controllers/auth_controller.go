package controller

import (
	"employee-attendance-system/internal/entity/dto"
	"employee-attendance-system/internal/middleware"
	"employee-attendance-system/internal/usecase"
	utils "employee-attendance-system/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthController interface {
	Signup(c *fiber.Ctx) error
	Signin(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	ChangeRole(c *fiber.Ctx) error
	Signout(c *fiber.Ctx) error
	UpdateProfile(c *fiber.Ctx) error
	GetProfile(ctx *fiber.Ctx) error
}

type authController struct {
	usecase  usecase.AuthUseCase
	log      *logrus.Logger
	validate *validator.Validate
}

func NewAuthController(usecase usecase.AuthUseCase, log *logrus.Logger, validate *validator.Validate) AuthController {
	return &authController{usecase: usecase, log: log, validate: validate}
}

func (c *authController) Signup(ctx *fiber.Ctx) error {
	var req dto.SignupRequest
	allowedFields := utils.GenerateAllowedFields(dto.SignupRequest{})
	if err := utils.BindAndValidateBody(ctx, &req, allowedFields, c.validate); err != nil {
		var errors []utils.ErrorDetail
		if validationErr := c.validate.Struct(req); validationErr != nil {
			for _, e := range validationErr.(validator.ValidationErrors) {
				errors = append(errors, utils.ErrorDetail{
					Field:   e.Field(),
					Message: e.Error(),
				})
			}
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), errors))
	}

	user, err := c.usecase.Signup(ctx.Context(), req.Email, req.Password, req.FullName)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessResponse(fiber.StatusCreated, "User created successfully", fiber.Map{
		"id":    user.ID,
		"email": user.Email,
	}, nil))
}

func (c *authController) Signin(ctx *fiber.Ctx) error {
	var req dto.SigninRequest
	allowedFields := utils.GenerateAllowedFields(dto.SigninRequest{})
	if err := utils.BindAndValidateBody(ctx, &req, allowedFields, c.validate); err != nil {
		var errors []utils.ErrorDetail
		if validationErr := c.validate.Struct(req); validationErr != nil {
			for _, e := range validationErr.(validator.ValidationErrors) {
				errors = append(errors, utils.ErrorDetail{
					Field:   e.Field(),
					Message: e.Error(),
				})
			}
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), errors))
	}
	deviceID := ctx.Get("X-Device-ID")
	if deviceID == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(utils.ErrorResponse(
			fiber.StatusUnprocessableEntity,
			"Validation failed",
			[]utils.ErrorDetail{{Field: "X-Device-ID", Message: "Device ID required"}},
		))
	}
	accessToken, refreshToken, user, err := c.usecase.Signin(ctx.Context(), req.Email, req.Password, &deviceID)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse(fiber.StatusUnauthorized, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Login successful", fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	}, nil))
}

func (c *authController) ChangePassword(ctx *fiber.Ctx) error {
	var req dto.ChangePasswordRequest
	allowedFields := utils.GenerateAllowedFields(dto.ChangePasswordRequest{})
	if err := utils.BindAndValidateBody(ctx, &req, allowedFields, c.validate); err != nil {
		var errors []utils.ErrorDetail
		if validationErr := c.validate.Struct(req); validationErr != nil {
			for _, e := range validationErr.(validator.ValidationErrors) {
				errors = append(errors, utils.ErrorDetail{
					Field:   e.Field(),
					Message: e.Error(),
				})
			}
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), errors))
	}

	localKeys := middleware.GetLocalKeys(ctx)
	if localKeys == nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "Internal Server error", nil))
	}
	if err := c.usecase.ChangePassword(ctx.Context(), localKeys.UserID, req.OldPassword, req.NewPassword); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Password changed successfully", nil, nil))
}

func (c *authController) RefreshToken(ctx *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	allowedFields := utils.GenerateAllowedFields(dto.RefreshTokenRequest{})
	if err := utils.BindAndValidateBody(ctx, &req, allowedFields, c.validate); err != nil {
		var errors []utils.ErrorDetail
		if validationErr := c.validate.Struct(req); validationErr != nil {
			for _, e := range validationErr.(validator.ValidationErrors) {
				errors = append(errors, utils.ErrorDetail{
					Field:   e.Field(),
					Message: e.Error(),
				})
			}
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), errors))
	}

	deviceID := ctx.Get("X-Device-ID")
	if deviceID == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(utils.ErrorResponse(
			fiber.StatusUnprocessableEntity,
			"Validation failed",
			[]utils.ErrorDetail{{Field: "X-Device-ID", Message: "Device ID required"}},
		))
	}

	newAccessToken, newRefreshToken, err := c.usecase.RefreshToken(ctx.Context(), req.RefreshToken, deviceID)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse(fiber.StatusUnauthorized, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Token refreshed successfully", fiber.Map{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken, // optional, bisa ganti token lama
	}, nil))
}

func (c *authController) ChangeRole(ctx *fiber.Ctx) error {
	var req dto.ChangeRoleRequest
	allowedFields := utils.GenerateAllowedFields(dto.ChangeRoleRequest{})
	if err := utils.BindAndValidateBody(ctx, &req, allowedFields, c.validate); err != nil {
		var errors []utils.ErrorDetail
		if validationErr := c.validate.Struct(req); validationErr != nil {
			for _, e := range validationErr.(validator.ValidationErrors) {
				errors = append(errors, utils.ErrorDetail{
					Field:   e.Field(),
					Message: e.Error(),
				})
			}
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), errors))
	}

	localKeys := middleware.GetLocalKeys(ctx)
	if localKeys == nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "Internal Server error", nil))
	}
	if err := c.usecase.ChangeRole(ctx.Context(), localKeys.UserID, req.Role); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Role changed successfully", nil, nil))
}

func (c *authController) Signout(ctx *fiber.Ctx) error {
	tokenHash := ctx.Get("Authorization") // Simulasi, ganti dengan ekstrak dari header
	if err := c.usecase.Signout(ctx.Context(), tokenHash); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Signout successful", nil, nil))
}

func (c *authController) UpdateProfile(ctx *fiber.Ctx) error {
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

func (c *authController) GetProfile(ctx *fiber.Ctx) error {
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
