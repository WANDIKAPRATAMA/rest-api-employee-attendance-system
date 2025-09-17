package dto

import (
	"employee-attendance-system/internal/entity/domain"
	"time"

	"github.com/google/uuid"
)

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
}

type SigninRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ChangeRoleRequest struct {
	UserID *uuid.UUID `json:"user_id" validate:"omitempty,uuid"`
	Role   string     `json:"role" validate:"required,oneof=employee admin "`
}

type UpdateProfileRequest struct {
	FullName  string `json:"full_name" validate:"omitempty,min=2,max=255"`
	Phone     string `json:"phone" validate:"omitempty,phone"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
	Address   string `json:"address" validate:"omitempty,max=500"`
}
type ProfileResponse struct {
	ID           uuid.UUID `json:"id"`
	SourceUserID uuid.UUID `json:"source_user_id"`
	EmployeeCode string    `json:"employee_code"`
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone"`
	AvatarURL    string    `json:"avatar_url"`
	Address      string    `json:"address"`
}

// Request untuk filter dynamic
type ListUsersRequest struct {
	Email          string     `query:"email" validate:"omitempty,email"`
	Status         string     `query:"status" validate:"omitempty,oneof=active inactive"`
	DepartmentID   *uuid.UUID `query:"department_id" validate:"omitempty,uuid"`
	CreatedAtStart *time.Time `query:"created_at_start" validate:"omitempty"`
	CreatedAtEnd   *time.Time `query:"created_at_end" validate:"omitempty"`

	Page  int `query:"page" validate:"omitempty,min=1"`          // Default 1
	Limit int `query:"limit" validate:"omitempty,min=1,max=100"` // Default 10
}

type UserResponse struct {
	ID              uuid.UUID               `json:"id"`
	SourceUserID    uuid.UUID               `json:"source_user_id"`
	EmployeeCode    string                  `json:"employee_code"`
	DepartmentID    *uuid.UUID              `json:"department_id,omitempty"`
	FullName        string                  `json:"full_name"`
	Phone           string                  `json:"phone"`
	AvatarURL       string                  `json:"avatar_url"`
	Address         string                  `json:"address"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
	Department      *DepartmentResponse     `json:"department,omitempty"`
	ApplicationRole *domain.ApplicationRole `json:"application_role,omitempty"`
	Email           string                  `json:"email"`
}
