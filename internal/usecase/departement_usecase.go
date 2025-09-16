// department_usecase.go
package usecase

import (
	"context"
	"employee-attendance-system/internal/entity/domain"
	"employee-attendance-system/internal/entity/dto"
	"employee-attendance-system/internal/repository"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DepartmentUseCase interface {
	CreateDepartment(ctx context.Context, req dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error)
	GetDepartment(ctx context.Context, id uuid.UUID) (*dto.DepartmentResponse, error)
	UpdateDepartment(ctx context.Context, id uuid.UUID, req dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error)
	DeleteDepartment(ctx context.Context, id uuid.UUID) error
	GetDepartments(ctx context.Context, page, limit int) ([]*dto.DepartmentResponse, int64, error)
	AssignmentDepartement(ctx context.Context, req dto.AssignmentDepartementRequest) error
}

type departmentUseCase struct {
	repo     repository.DepartmentRepository
	userRepo repository.UserRepository
	log      *logrus.Logger
	validate *validator.Validate
}

func NewDepartmentUseCase(repo repository.DepartmentRepository, log *logrus.Logger, validate *validator.Validate, userRepo repository.UserRepository) DepartmentUseCase {
	return &departmentUseCase{repo: repo, log: log, validate: validate, userRepo: userRepo}
}

func (u *departmentUseCase) AssignmentDepartement(ctx context.Context, req dto.AssignmentDepartementRequest) error {
	if exist, err := u.userRepo.IsUserExist(req.UserID); !exist || err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user not found")
		}
		return err
	}

	if role, err := u.userRepo.FindUserRoleByUserID(req.UserID); role == domain.Admin || err != nil {
		if role == domain.Admin {
			return fmt.Errorf("admin cannot be assigned to department")
		}
		return err
	}

	if exist, err := u.repo.IsDepartmentExist(req.DepartmentID); !exist || err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("department not found")
		}
		return err
	}

	if err := u.repo.AssignmentDepartement(req.UserID, req.DepartmentID); err != nil {
		return err
	}

	return nil
}
func (u *departmentUseCase) CreateDepartment(ctx context.Context, req dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error) {
	// Parse only time
	layout := "15:04:05"
	clockIn, err := time.Parse(layout, req.MaxClockInTime)
	if err != nil {
		return nil, fmt.Errorf("invalid max_clock_in_time: %w", err)
	}
	clockOut, err := time.Parse(layout, req.MaxClockOutTime)
	if err != nil {
		return nil, fmt.Errorf("invalid max_clock_out_time: %w", err)
	}

	dept := &domain.Department{
		Name:            req.Name,
		MaxClockInTime:  clockIn,
		MaxClockOutTime: clockOut,
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

	// kalau user mengisi clock in time
	if !req.MaxClockInTime.IsZero() {
		// normalisasi biar hanya jam (hilangkan tanggal)
		dept.MaxClockInTime = time.Date(0, 1, 1,
			req.MaxClockInTime.Hour(),
			req.MaxClockInTime.Minute(),
			req.MaxClockInTime.Second(),
			0,
			time.UTC,
		)
	}

	// kalau user mengisi clock out time
	if !req.MaxClockOutTime.IsZero() {
		// normalisasi biar hanya jam (hilangkan tanggal)
		dept.MaxClockOutTime = time.Date(0, 1, 1,
			req.MaxClockOutTime.Hour(),
			req.MaxClockOutTime.Minute(),
			req.MaxClockOutTime.Second(),
			0,
			time.UTC,
		)
	}

	if err := u.repo.UpdateDepartment(dept); err != nil {
		return nil, err
	}

	return mapToDepartmentResponse(dept), nil
}

func (u *departmentUseCase) DeleteDepartment(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteDepartment(id)
}

/*************  ✨ Windsurf Command ⭐  *************/
/*******  2e04a59a-ec79-4895-bb60-83859705510f  *******/
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
	if d == nil {
		return nil
	}
	return &dto.DepartmentResponse{
		ID:              d.ID,
		Name:            d.Name,
		MaxClockInTime:  d.MaxClockInTime,
		MaxClockOutTime: d.MaxClockOutTime,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}
