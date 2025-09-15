package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// New struct for Department
type Department struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name            string         `gorm:"column:department_name;type:varchar(255);not null"`
	MaxClockInTime  time.Time      `gorm:"type:time;not null"` // Time only (e.g., 09:00:00)
	MaxClockOutTime time.Time      `gorm:"type:time;not null"` // Time only (e.g., 17:00:00)
	CreatedAt       time.Time      `gorm:"default:current_timestamp"`
	UpdatedAt       time.Time      `gorm:"default:current_timestamp"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

// New struct for Attendance (daily record)
type Attendance struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EmployeeCode string         `gorm:"type:varchar(50);index;not null"` // FK to UserProfile.EmployeeCode
	AttendanceID string         `gorm:"type:varchar(100);uniqueIndex"`   // Unique ID (e.g., generate as "EMP001-2025-09-15")
	ClockIn      *time.Time     `gorm:"type:timestamp"`                  // Nullable
	ClockOut     *time.Time     `gorm:"type:timestamp"`                  // Nullable
	CreatedAt    time.Time      `gorm:"default:current_timestamp"`
	UpdatedAt    time.Time      `gorm:"default:current_timestamp"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// New struct for AttendanceHistory (logs for each in/out action)

type AttendanceType string

const (
	AttendanceTypeIn  AttendanceType = "in"
	AttendanceTypeOut AttendanceType = "out"
)

type AttendanceHistory struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EmployeeCode   string         `gorm:"type:varchar(50);index;not null"`  // FK to UserProfile.EmployeeCode
	AttendanceID   string         `gorm:"type:varchar(100);index;not null"` // FK to Attendance.AttendanceID
	DateAttendance time.Time      `gorm:"type:timestamp;not null"`
	AttendanceType AttendanceType `gorm:"type:attendance_type;not null"` // in = In, out = Out
	Description    string         `gorm:"type:text"`
	CreatedAt      time.Time      `gorm:"default:current_timestamp"`
	UpdatedAt      time.Time      `gorm:"default:current_timestamp"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
