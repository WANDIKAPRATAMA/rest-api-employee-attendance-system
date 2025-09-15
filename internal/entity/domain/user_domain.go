package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Existing structs (unchanged)
type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email         string         `gorm:"type:varchar(255);unique;not null"`
	Status        string         `gorm:"type:user_status;not null;default:'inactive'"`
	EmailVerified bool           `gorm:"column:email_verified;not null;default:false"`
	CreatedAt     time.Time      `gorm:"default:current_timestamp"`
	UpdatedAt     time.Time      `gorm:"default:current_timestamp"`
	DeletedAt     gorm.DeletedAt `gorm:"index"` // Soft delete
}

type UserSecurity struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SourceUserID uuid.UUID      `gorm:"column:source_user_id;type:uuid;not null"`
	Password     string         `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time      `gorm:"default:current_timestamp"`
	UpdatedAt    time.Time      `gorm:"default:current_timestamp"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// Updated UserProfile (integrated with Employee details)
type UserProfile struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SourceUserID uuid.UUID      `gorm:"column:source_user_id;type:uuid;not null;uniqueIndex"` // Links to User.ID
	EmployeeCode string         `gorm:"column:employee_code;type:varchar(50);uniqueIndex"`    // Matches employee.employee_id (unique code)
	DepartmentID *uuid.UUID     `gorm:"type:uuid;index"`                                      // FK to Department.ID (nullable if not assigned)
	FullName     string         `gorm:"type:varchar(255)"`                                    // Matches employee.name
	Phone        string         `gorm:"type:varchar(50)"`
	AvatarURL    string         `gorm:"type:text"`
	Address      string         `gorm:"type:text"` // Matches employee.address
	CreatedAt    time.Time      `gorm:"default:current_timestamp"`
	UpdatedAt    time.Time      `gorm:"default:current_timestamp"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	// Relationships (optional for preloading)
	Department Department `gorm:"foreignKey:DepartmentID"`
}
