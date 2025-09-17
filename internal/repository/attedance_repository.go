// attendance_repository.go
package repository

import (
	"employee-attendance-system/internal/entity/domain"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AttendanceRepository interface {
	FindAttendanceByID(attendanceID string, attendance *domain.Attendance) error
	CreateAttendanceWithHistory(attendance *domain.Attendance, history *domain.AttendanceHistory) error
	UpdateAttendanceWithHistory(attendance *domain.Attendance, history *domain.AttendanceHistory) error
	GetAttendanceQuery() *gorm.DB
	FindAttendanceHistoryByEmployeeCode(employeeCode string, page, limit int) ([]*domain.AttendanceHistory, int64, error)

	FindCurrentAttendance(employeeCode string) (*domain.Attendance, error)
}

type attendanceRepository struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewAttendanceRepository(db *gorm.DB, log *logrus.Logger) AttendanceRepository {
	return &attendanceRepository{db: db, log: log}
}

func (r *attendanceRepository) FindAttendanceByID(attendanceID string, attendance *domain.Attendance) error {
	return r.db.Where("attendance_id = ?", attendanceID).First(attendance).Error
}

func (r *attendanceRepository) CreateAttendanceWithHistory(attendance *domain.Attendance, history *domain.AttendanceHistory) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(attendance).Error; err != nil {
			return err
		}
		return tx.Create(history).Error
	})
}

func (r *attendanceRepository) UpdateAttendanceWithHistory(attendance *domain.Attendance, history *domain.AttendanceHistory) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(attendance).Error; err != nil {
			return err
		}
		return tx.Create(history).Error
	})
}
func (r *attendanceRepository) FindAttendanceHistoryByEmployeeCode(employeeCode string, page, limit int) ([]*domain.AttendanceHistory, int64, error) {
	var histories []*domain.AttendanceHistory
	query := r.db.Model(&domain.AttendanceHistory{}).
		Where("employee_code = ? AND deleted_at IS NULL", employeeCode)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("date_attendance DESC").Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}
func (r *attendanceRepository) GetAttendanceQuery() *gorm.DB {
	return r.db.Table("attendances a").
		Joins("JOIN user_profiles up ON a.employee_code = up.employee_code").
		Joins("JOIN departments d ON up.department_id = d.id").
		Select(`
			a.attendance_id,
			a.employee_code,
			up.full_name,
			d.department_name,
			a.clock_in,
			a.clock_out,
			d.max_clock_in_time,
			d.max_clock_out_time
		`)
}

func (r *attendanceRepository) FindCurrentAttendance(employeeCode string) (*domain.Attendance, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()) // 10:06 AM WIB, 17 Sept 2025

	var attendance domain.Attendance
	err := r.db.Where("employee_code = ? AND DATE(created_at) = ? AND deleted_at IS NULL", employeeCode, today).
		Order("created_at DESC").First(&attendance).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &attendance, nil
}
