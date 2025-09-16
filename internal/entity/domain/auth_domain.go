package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	Employee Role = "employee"
	Admin    Role = "admin"
)

type ApplicationRole struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	SourceUserID uuid.UUID      `gorm:"column:source_user_id;type:uuid;not null" json:"source_user_id"`
	Role         Role           `gorm:"type:app_role;not null" json:"role"`
	CreatedAt    time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type RefreshToken struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SourceUserID uuid.UUID `gorm:"column:source_user_id;type:uuid;not null;uniqueIndex:idx_user_device"`
	DeviceID     string    `gorm:"type:text;not null;uniqueIndex:idx_user_device"`
	TokenHash    string    `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"default:current_timestamp"`
	ExpiresAt    time.Time `gorm:"not null"`
	LastUsedAt   time.Time
	RevokedAt    *time.Time     `gorm:"column:revoked_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
