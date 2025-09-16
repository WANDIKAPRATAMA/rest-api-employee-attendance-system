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

// func (u *attendanceUseCase) GetAttendanceLogs(ctx context.Context, userID uuid.UUID, role string, req dto.GetAttendanceLogsRequest) ([]dto.AttendanceLogResponse, int64, error) {
// 	offset := (req.Page - 1) * req.Limit

// 	// Dynamic query
// 	query := u.repo.GetAttendanceQuery()

// 	if req.Date != "" {
// 		query = query.Where("DATE(a.clock_in) = ?", req.Date)
// 	}
// 	if req.DepartmentID != nil {
// 		query = query.Where("up.department_id = ?", *req.DepartmentID)
// 	}

// 	// Jika bukan admin, limit ke own department
// 	if role != "admin" {
// 		profile, _ := u.profileRepo.FindUserProfileByUserID(userID)
// 		if profile != nil && profile.DepartmentID != nil {
// 			query = query.Where("up.department_id = ?", *profile.DepartmentID)
// 		} else {
// 			return nil, 0, fmt.Errorf("no access")
// 		}
// 	}

// 	var total int64
// 	if err := query.Count(&total).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	var logs []dto.AttendanceLogResponse
// 	if err := query.Offset(offset).Limit(req.Limit).Scan(&logs).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	return logs, total, nil
// }

func (u *attendanceUseCase) GetAttendanceLogs(ctx context.Context, userID uuid.UUID, role string, req dto.GetAttendanceLogsRequest) ([]dto.AttendanceLogResponse, int64, error) {
	u.log.WithFields(logrus.Fields{
		"user_id":       userID,
		"role":          role,
		"page":          req.Page,
		"limit":         req.Limit,
		"date":          req.Date,
		"department_id": req.DepartmentID,
	}).Info("Starting GetAttendanceLogs")

	offset := (req.Page - 1) * req.Limit

	query := u.repo.GetAttendanceQuery()

	if req.Date != "" {
		u.log.WithField("filter_date", req.Date).Debug("Applying date filter")
		query = query.Where("DATE(a.clock_in) = ?", req.Date)
	}
	if req.DepartmentID != nil {
		u.log.WithField("filter_department_id", req.DepartmentID).Debug("Applying department filter")
		query = query.Where("up.department_id = ?", *req.DepartmentID)
	}

	if role != "admin" {
		profile, _ := u.profileRepo.FindUserProfileByUserID(userID)
		u.log.WithField("user_profile", profile).Debug("Non-admin role, checking department restriction")
		if profile != nil && profile.DepartmentID != nil {
			query = query.Where("up.department_id = ?", *profile.DepartmentID)
		} else {
			u.log.Warn("No access: user has no department")
			return nil, 0, fmt.Errorf("no access")
		}
	}

	var total int64
	if err := query.Model(&dto.RawAttendanceLog{}).Count(&total).Error; err != nil {
		u.log.WithError(err).Error("Failed to count attendance logs")
		return nil, 0, err
	}
	u.log.WithField("total_records", total).Info("Total attendance logs found")

	var rawLogs []dto.RawAttendanceLog
	if err := query.Offset(offset).Limit(req.Limit).Scan(&rawLogs).Error; err != nil {
		u.log.WithError(err).Error("Failed to scan raw attendance logs")
		return nil, 0, err
	}
	u.log.WithField("raw_logs", rawLogs).Debug("Fetched raw attendance logs from DB")

	finalLogs := make([]dto.AttendanceLogResponse, 0, len(rawLogs))

	for _, raw := range rawLogs {
		inPunctuality := "N/A"
		outPunctuality := "N/A"

		u.log.WithFields(logrus.Fields{
			"attendance_id":      raw.AttendanceID,
			"employee_code":      raw.EmployeeCode,
			"clock_in":           raw.ClockIn,
			"max_clock_in_time":  raw.MaxClockInTime,
			"clock_out":          raw.ClockOut,
			"max_clock_out_time": raw.MaxClockOutTime,
		}).Debug("Processing attendance record")

		// Kalkulasi Punctuality Clock In
		if raw.ClockIn != nil {
			if raw.MaxClockInTime == nil {
				u.log.WithField("department", raw.DepartmentName).Error("MaxClockInTime not configured")
				return nil, 0, fmt.Errorf("konfigurasi 'MaxClockInTime' untuk departemen '%s' tidak ditemukan", raw.DepartmentName)
			}

			maxIn := *raw.MaxClockInTime
			actualClockIn := *raw.ClockIn
			targetInTime := time.Date(
				actualClockIn.Year(), actualClockIn.Month(), actualClockIn.Day(),
				maxIn.Hour(), maxIn.Minute(), maxIn.Second(), 0, actualClockIn.Location(),
			)

			u.log.WithFields(logrus.Fields{
				"actual_clock_in": actualClockIn,
				"target_in_time":  targetInTime,
			}).Debug("Comparing clock in times")

			if actualClockIn.After(targetInTime) {
				inPunctuality = "Late"
			} else {
				inPunctuality = "On Time"
			}
		}

		// Kalkulasi Punctuality Clock Out
		// Kalkulasi Punctuality Clock Out
		if raw.ClockOut != nil {
			if raw.MaxClockOutTime == nil {
				return nil, 0, fmt.Errorf("konfigurasi 'MaxClockOutTime' untuk departemen '%s' tidak ditemukan", raw.DepartmentName)
			}

			actualClockOut := *raw.ClockOut
			maxOut := *raw.MaxClockOutTime

			// Gunakan tanggal actualClockOut, jam dari MaxClockOutTime
			targetOutTime := time.Date(
				actualClockOut.Year(), actualClockOut.Month(), actualClockOut.Day(),
				maxOut.Hour(), maxOut.Minute(), maxOut.Second(), 0, actualClockOut.Location(),
			)

			u.log.WithFields(logrus.Fields{
				"actual_clock_out": actualClockOut,
				"target_out_time":  targetOutTime,
			}).Debug("Comparing clock out times")

			if actualClockOut.Before(targetOutTime) {
				outPunctuality = "Early Leave"
			} else {
				outPunctuality = "On Time"
			}
		}

		u.log.WithFields(logrus.Fields{
			"in_punctuality":  inPunctuality,
			"out_punctuality": outPunctuality,
		}).Info("Calculated punctuality result")

		finalLogs = append(finalLogs, dto.AttendanceLogResponse{
			AttendanceID:   raw.AttendanceID,
			EmployeeCode:   raw.EmployeeCode,
			FullName:       raw.FullName,
			DepartmentName: raw.DepartmentName,
			ClockIn:        raw.ClockIn,
			ClockOut:       raw.ClockOut,
			InPunctuality:  inPunctuality,
			OutPunctuality: outPunctuality,
		})
	}

	return finalLogs, total, nil
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
func parseDepartmentTime(base *time.Time, deptTime *time.Time, fallback string) time.Time {
	var t time.Time

	if deptTime != nil {
		t = *deptTime
	} else {
		parsed, _ := time.Parse("15:04:05", fallback)
		t = parsed
	}

	// gabungkan tanggal dari base (clock_in/clock_out) dengan jam aturan dept
	year, month, day := base.Date()
	return time.Date(year, month, day, t.Hour(), t.Minute(), t.Second(), 0, base.Location())
}
