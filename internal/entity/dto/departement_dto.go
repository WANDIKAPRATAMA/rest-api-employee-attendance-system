package dto

import (
	"time"

	"github.com/google/uuid"
)

// Untuk Department
type CreateDepartmentRequest struct {
	Name            string    `json:"name" validate:"required,min=3,max=255"`
	MaxClockInTime  time.Time `json:"max_clock_in_time" validate:"required"`  // e.g., "09:00:00"
	MaxClockOutTime time.Time `json:"max_clock_out_time" validate:"required"` // e.g., "17:00:00"
}

type UpdateDepartmentRequest struct {
	Name            string    `json:"name" validate:"omitempty,min=3,max=255"`
	MaxClockInTime  time.Time `json:"max_clock_in_time" validate:"omitempty"`
	MaxClockOutTime time.Time `json:"max_clock_out_time" validate:"omitempty"`
}

type DepartmentResponse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	MaxClockInTime  time.Time `json:"max_clock_in_time"`
	MaxClockOutTime time.Time `json:"max_clock_out_time"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Untuk Attendance
type ClockInRequest struct {
	// Kosong, karena auto-detect dari user
}

type ClockOutRequest struct {
	// Kosong, sama
}

type AttendanceResponse struct {
	ID           uuid.UUID  `json:"id"`
	EmployeeCode string     `json:"employee_code"`
	AttendanceID string     `json:"attendance_id"`
	ClockIn      *time.Time `json:"clock_in"`
	ClockOut     *time.Time `json:"clock_out"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type AttendanceLogResponse struct {
	AttendanceID   string     `json:"attendance_id"`
	EmployeeCode   string     `json:"employee_code"`
	FullName       string     `json:"full_name"`
	DepartmentName string     `json:"department_name"`
	ClockIn        *time.Time `json:"clock_in"`
	ClockOut       *time.Time `json:"clock_out"`
	InPunctuality  string     `json:"in_punctuality"`  // "On Time" or "Late"
	OutPunctuality string     `json:"out_punctuality"` // "On Time" or "Early Leave"
}

// Untuk filters di GET logs
type GetAttendanceLogsRequest struct {
	Date         string     `query:"date" validate:"omitempty,datetime=2006-01-02"` // YYYY-MM-DD
	DepartmentID *uuid.UUID `query:"department_id" validate:"omitempty,uuid"`
	Page         int        `query:"page" validate:"omitempty,min=1"`          // Default 1
	Limit        int        `query:"limit" validate:"omitempty,min=1,max=100"` // Default 10
}
