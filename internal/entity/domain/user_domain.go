package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Email         string         `gorm:"type:varchar(255);unique;not null" json:"email"`
	Status        string         `gorm:"type:user_status;not null;default:'inactive'" json:"status"`
	EmailVerified bool           `gorm:"column:email_verified;not null;default:false" json:"email_verified"`
	CreatedAt     time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"` // Soft delete
}

type UserSecurity struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	SourceUserID uuid.UUID      `gorm:"column:source_user_id;type:uuid;not null" json:"source_user_id"`
	Password     string         `gorm:"type:varchar(255);not null" json:"password"`
	CreatedAt    time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type UserProfile struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	SourceUserID uuid.UUID      `gorm:"column:source_user_id;type:uuid;not null;uniqueIndex" json:"source_user_id"`
	EmployeeCode string         `gorm:"column:employee_code;type:varchar(50);uniqueIndex" json:"employee_code"`
	DepartmentID *uuid.UUID     `gorm:"type:uuid;index" json:"department_id,omitempty"`
	FullName     string         `gorm:"type:varchar(255)" json:"full_name"`
	Phone        string         `gorm:"type:varchar(50)" json:"phone"`
	AvatarURL    string         `gorm:"type:text" json:"avatar_url"`
	Address      string         `gorm:"type:text" json:"address"`
	CreatedAt    time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Department      *Department      `gorm:"foreignKey:DepartmentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"department,omitempty"`
	ApplicationRole *ApplicationRole `gorm:"foreignKey:SourceUserID;references:SourceUserID" json:"application_role,omitempty"`
}
