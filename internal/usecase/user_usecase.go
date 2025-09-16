package usecase

import (
	"context"
	"employee-attendance-system/internal/entity/domain"
	"employee-attendance-system/internal/entity/dto"
	"employee-attendance-system/internal/repository"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserUseCase interface {
	ListUsers(ctx context.Context, req dto.ListUsersRequest) ([]*dto.UserResponse, int64, error)

	UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest) (*domain.UserProfile, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserProfile, error)
}

type userUseCase struct {
	repo     repository.UserRepository
	log      *logrus.Logger
	validate *validator.Validate
}

func NewUserUseCase(repo repository.UserRepository, log *logrus.Logger, validate *validator.Validate) UserUseCase {
	return &userUseCase{repo: repo, log: log, validate: validate}
}

func mapToUserResponse(up *domain.UserProfile) *dto.UserResponse {

	return &dto.UserResponse{
		ID:              up.ID,
		SourceUserID:    up.SourceUserID,
		EmployeeCode:    up.EmployeeCode,
		DepartmentID:    up.DepartmentID,
		FullName:        up.FullName,
		Phone:           up.Phone,
		AvatarURL:       up.AvatarURL,
		Address:         up.Address,
		CreatedAt:       up.CreatedAt,
		UpdatedAt:       up.UpdatedAt,
		Department:      mapToDepartmentResponse(up.Department),
		ApplicationRole: up.ApplicationRole,
	}
}

func (u *userUseCase) UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest) (*domain.UserProfile, error) {
	// Cari profile existing
	profile, err := u.repo.FindUserProfileByUserID(userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, fmt.Errorf("profile not found")
	}

	// Update fields (hanya yang diisi)
	if req.FullName != "" {
		profile.FullName = req.FullName
	}
	if req.Phone != "" {
		profile.Phone = req.Phone
	}
	if req.AvatarURL != "" {
		profile.AvatarURL = req.AvatarURL
	}
	if req.Address != "" {
		profile.Address = req.Address
	}

	if err := u.repo.UpdateUserProfile(profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (u *userUseCase) ListUsers(ctx context.Context, req dto.ListUsersRequest) ([]*dto.UserResponse, int64, error) {
	// Build dynamic query logic
	users, total, err := u.repo.FindAllUsers(req)
	if err != nil {
		return nil, 0, err
	}

	res := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		res[i] = mapToUserResponse(user)
	}

	return res, total, nil
}

func (u *userUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserProfile, error) {
	profile, err := u.repo.FindUserProfileByUserID(userID)
	if err != nil {
		return nil, err
	}
	if profile.DepartmentID == nil {
		profile.Department = nil
	}
	return profile, nil
}
