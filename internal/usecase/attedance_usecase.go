// attendance_usecase.go
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

type AttendanceUseCase interface {
	ClockIn(ctx context.Context, userID uuid.UUID) (*dto.AttendanceResponse, error)
	ClockOut(ctx context.Context, userID uuid.UUID) (*dto.AttendanceResponse, error)
	GetAttendanceLogs(ctx context.Context, userID uuid.UUID, role string, req dto.GetAttendanceLogsRequest) ([]dto.AttendanceLogResponse, int64, error)
}

type attendanceUseCase struct {
	repo        repository.AttendanceRepository
	profileRepo repository.UserRepository // Untuk get employee code
	log         *logrus.Logger
	validate    *validator.Validate
}

func NewAttendanceUseCase(repo repository.AttendanceRepository, profileRepo repository.UserRepository, log *logrus.Logger, validate *validator.Validate) AttendanceUseCase {
	return &attendanceUseCase{repo: repo, profileRepo: profileRepo, log: log, validate: validate}

}
func (u *attendanceUseCase) ClockIn(ctx context.Context, userID uuid.UUID) (*dto.AttendanceResponse, error) {
	profile, err := u.profileRepo.FindUserProfileByUserID(userID)
	if err != nil || profile == nil {
		return nil, fmt.Errorf("profile not found")
	}
	if profile.DepartmentID == nil {
		return nil, fmt.Errorf("no department assigned")
	}

	now := time.Now()
	today := now.Format("2006-01-02")
	attendanceID := fmt.Sprintf("%s-%s", profile.EmployeeCode, today)

	var attendance domain.Attendance
	if err := u.repo.FindAttendanceByID(attendanceID, &attendance); err == nil {
		return nil, fmt.Errorf("already clocked in today")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	attendance = domain.Attendance{
		EmployeeCode: profile.EmployeeCode,
		AttendanceID: attendanceID,
		ClockIn:      &now,
	}

	history := domain.AttendanceHistory{
		EmployeeCode:   profile.EmployeeCode,
		AttendanceID:   attendanceID,
		DateAttendance: now,
		AttendanceType: domain.AttendanceTypeIn,
		Description:    "Clock in",
	}

	err = u.repo.CreateAttendanceWithHistory(&attendance, &history)
	if err != nil {
		return nil, err
	}

	return mapToAttendanceResponse(&attendance), nil
}

func (u *attendanceUseCase) ClockOut(ctx context.Context, userID uuid.UUID) (*dto.AttendanceResponse, error) {
	profile, err := u.profileRepo.FindUserProfileByUserID(userID)
	if err != nil || profile == nil {
		return nil, fmt.Errorf("profile not found")
	}

	now := time.Now()
	today := now.Format("2006-01-02")
	attendanceID := fmt.Sprintf("%s-%s", profile.EmployeeCode, today)

	var attendance domain.Attendance
	if err := u.repo.FindAttendanceByID(attendanceID, &attendance); err != nil {
		return nil, fmt.Errorf("no clock in today")
	}
	if attendance.ClockOut != nil {
		return nil, fmt.Errorf("already clocked out")
	}

	attendance.ClockOut = &now

	history := domain.AttendanceHistory{
		EmployeeCode:   profile.EmployeeCode,
		AttendanceID:   attendanceID,
		DateAttendance: now,
		AttendanceType: domain.AttendanceTypeOut,
		Description:    "Clock out",
	}

	err = u.repo.UpdateAttendanceWithHistory(&attendance, &history)
	if err != nil {
		return nil, err
	}

	return mapToAttendanceResponse(&attendance), nil
}

func (u *attendanceUseCase) GetAttendanceLogs(ctx context.Context, userID uuid.UUID, role string, req dto.GetAttendanceLogsRequest) ([]dto.AttendanceLogResponse, int64, error) {
	offset := (req.Page - 1) * req.Limit

	// Dynamic query
	query := u.repo.GetAttendanceQuery()

	if req.Date != "" {
		query = query.Where("DATE(a.clock_in) = ?", req.Date)
	}
	if req.DepartmentID != nil {
		query = query.Where("up.department_id = ?", *req.DepartmentID)
	}

	// Jika bukan admin, limit ke own department
	if role != "admin" {
		profile, _ := u.profileRepo.FindUserProfileByUserID(userID)
		if profile != nil && profile.DepartmentID != nil {
			query = query.Where("up.department_id = ?", *profile.DepartmentID)
		} else {
			return nil, 0, fmt.Errorf("no access")
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []dto.AttendanceLogResponse
	if err := query.Offset(offset).Limit(req.Limit).Scan(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func mapToAttendanceResponse(a *domain.Attendance) *dto.AttendanceResponse {
	return &dto.AttendanceResponse{
		ID:           a.ID,
		EmployeeCode: a.EmployeeCode,
		AttendanceID: a.AttendanceID,
		ClockIn:      a.ClockIn,
		ClockOut:     a.ClockOut,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}
