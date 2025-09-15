// attendance_repository.go
package repository

import (
	"employee-attendance-system/internal/entity/domain"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AttendanceRepository interface {
	FindAttendanceByID(attendanceID string, attendance *domain.Attendance) error
	CreateAttendanceWithHistory(attendance *domain.Attendance, history *domain.AttendanceHistory) error
	UpdateAttendanceWithHistory(attendance *domain.Attendance, history *domain.AttendanceHistory) error
	GetAttendanceQuery() *gorm.DB // Untuk dynamic query
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

func (r *attendanceRepository) GetAttendanceQuery() *gorm.DB {
	return r.db.Table("attendances a").
		Joins("JOIN user_profiles up ON a.employee_code = up.employee_code").
		Joins("JOIN departments d ON up.department_id = d.id").
		Select(
			"a.attendance_id, a.employee_code, up.full_name, d.name as department_name, a.clock_in, a.clock_out, " +
				"CASE WHEN a.clock_in > (DATE(a.clock_in) + d.max_clock_in_time) THEN 'Late' ELSE 'On Time' END as in_punctuality, " +
				"CASE WHEN a.clock_out < (DATE(a.clock_out) + d.max_clock_out_time) THEN 'Early Leave' ELSE 'On Time' END as out_punctuality",
		)
}
