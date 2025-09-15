// department_usecase.go
package usecase

import (
	"context"
	"employee-attendance-system/internal/entity/domain"
	"employee-attendance-system/internal/entity/dto"
	"employee-attendance-system/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type DepartmentUseCase interface {
	CreateDepartment(ctx context.Context, req dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error)
	GetDepartment(ctx context.Context, id uuid.UUID) (*dto.DepartmentResponse, error)
	UpdateDepartment(ctx context.Context, id uuid.UUID, req dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error)
	DeleteDepartment(ctx context.Context, id uuid.UUID) error
	GetDepartments(ctx context.Context, page, limit int) ([]*dto.DepartmentResponse, int64, error)
}

type departmentUseCase struct {
	repo     repository.DepartmentRepository
	log      *logrus.Logger
	validate *validator.Validate
}

func NewDepartmentUseCase(repo repository.DepartmentRepository, log *logrus.Logger, validate *validator.Validate) DepartmentUseCase {
	return &departmentUseCase{repo: repo, log: log, validate: validate}
}

func (u *departmentUseCase) CreateDepartment(ctx context.Context, req dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error) {
	dept := &domain.Department{
		Name:            req.Name,
		MaxClockInTime:  req.MaxClockInTime,
		MaxClockOutTime: req.MaxClockOutTime,
	}
	if err := u.repo.CreateDepartment(dept); err != nil {
		return nil, err
	}
	return mapToDepartmentResponse(dept), nil
}

func (u *departmentUseCase) GetDepartment(ctx context.Context, id uuid.UUID) (*dto.DepartmentResponse, error) {
	dept, err := u.repo.FindDepartmentByID(id)
	if err != nil {
		return nil, err
	}
	return mapToDepartmentResponse(dept), nil
}

func (u *departmentUseCase) UpdateDepartment(ctx context.Context, id uuid.UUID, req dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error) {
	dept, err := u.repo.FindDepartmentByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		dept.Name = req.Name
	}
	if !req.MaxClockInTime.IsZero() {
		dept.MaxClockInTime = req.MaxClockInTime
	}
	if !req.MaxClockOutTime.IsZero() {
		dept.MaxClockOutTime = req.MaxClockOutTime
	}
	if err := u.repo.UpdateDepartment(dept); err != nil {
		return nil, err
	}
	return mapToDepartmentResponse(dept), nil
}

func (u *departmentUseCase) DeleteDepartment(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteDepartment(id)
}

func (u *departmentUseCase) GetDepartments(ctx context.Context, page, limit int) ([]*dto.DepartmentResponse, int64, error) {
	offset := (page - 1) * limit
	depts, total, err := u.repo.FindAllDepartments(offset, limit)
	if err != nil {
		return nil, 0, err
	}
	res := make([]*dto.DepartmentResponse, len(depts))
	for i, d := range depts {
		res[i] = mapToDepartmentResponse(d)
	}
	return res, total, nil
}

func mapToDepartmentResponse(d *domain.Department) *dto.DepartmentResponse {
	return &dto.DepartmentResponse{
		ID:              d.ID,
		Name:            d.Name,
		MaxClockInTime:  d.MaxClockInTime,
		MaxClockOutTime: d.MaxClockOutTime,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}
