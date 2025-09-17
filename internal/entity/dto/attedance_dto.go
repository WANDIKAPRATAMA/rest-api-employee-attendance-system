package dto

import (
	"time"

	"github.com/google/uuid"
)

type GetAttendanceHistoryRequest struct {
	UserID uuid.UUID `query:"user_id" validate:"required,uuid"`
	Page   int       `query:"page" validate:"omitempty,min=1"`          // Default 1
	Limit  int       `query:"limit" validate:"omitempty,min=1,max=100"` // Default 10
}

type AttendanceHistoryResponse struct {
	ID             uuid.UUID `json:"id"`
	EmployeeCode   string    `json:"employee_code"`
	AttendanceID   string    `json:"attendance_id"`
	DateAttendance time.Time `json:"date_attendance"`
	AttendanceType string    `json:"attendance_type"` // "in" or "out"
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Untuk Admin Dashboard
type AdminDashboardRequest struct {
	StartDate *time.Time `query:"start_date" validate:"omitempty,datetime"` // e.g., 2025-09-01
	EndDate   *time.Time `query:"end_date" validate:"omitempty,datetime"`
}

type AdminDashboardResponse struct {
	TotalEmployeesPerDept   map[string]int `json:"total_employees_per_dept"` // Key: Dept Name, Value: Count
	TotalUpdatedDepts       int            `json:"total_updated_depts"`
	TotalTodayRegistrations int            `json:"total_today_registrations"`
}

type CheckCurrentStatusRequest struct {
	UserID *uuid.UUID `query:"user_id" validate:"omitempty,uuid"` // Kosong berarti cek diri sendiri
}

// Response untuk status saat ini
type CurrentStatusResponse struct {
	UserID       uuid.UUID  `json:"user_id"`
	EmployeeCode string     `json:"employee_code"`
	FullName     string     `json:"full_name"`
	Department   string     `json:"department,omitempty"`
	Status       string     `json:"status"` // "Clocked In", "Clocked Out", "Not Clocked"
	ClockIn      *time.Time `json:"clock_in,omitempty"`
	ClockOut     *time.Time `json:"clock_out,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
