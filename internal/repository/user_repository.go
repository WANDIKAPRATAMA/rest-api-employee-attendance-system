package repository

import (
	"employee-attendance-system/internal/entity/domain"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	CreateUser(user *domain.User, profile *domain.UserProfile, security *domain.UserSecurity, role *domain.ApplicationRole) error
	FindUserByEmail(email string) (*domain.User, error)
	FindUserSecurityByUserID(userID uuid.UUID) (*domain.UserSecurity, error)
	CreateRefreshToken(token *domain.RefreshToken) error
	FindRefreshToken(token string, deviceID string) (*domain.RefreshToken, error)
	RevokeRefreshToken(tokenHash string) error
	UpdateUserSecurity(userID uuid.UUID, newPassword string) error
	AssignRole(userID uuid.UUID, role domain.Role) error
	FindUserRoleByUserID(userID uuid.UUID) (domain.Role, error)
	FindUserByID(user_id uuid.UUID) (*domain.User, error)
	UpdateRefreshToken(token *domain.RefreshToken) error
	UpdateUserProfile(profile *domain.UserProfile) error
	FindUserProfileByUserID(userID uuid.UUID) (*domain.UserProfile, error)
}

type userRepository struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewUserRepository(db *gorm.DB, log *logrus.Logger) UserRepository {
	return &userRepository{db: db, log: log}
}

func (r *userRepository) FindUserRoleByUserID(userID uuid.UUID) (domain.Role, error) {
	var role domain.ApplicationRole
	if err := r.db.Where("source_user_id = ?", userID).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Employee, nil
		}
		return "", err
	}
	return role.Role, nil
}

func (r *userRepository) CreateUser(user *domain.User, profile *domain.UserProfile, security *domain.UserSecurity, role *domain.ApplicationRole) error {
	tx := r.db.Begin()
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	profile.SourceUserID = user.ID
	security.SourceUserID = user.ID
	role.SourceUserID = user.ID
	if err := tx.Create(profile).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(security).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(role).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *userRepository) FindUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindUserByID(user_id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("id = ?", user_id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindUserSecurityByUserID(userID uuid.UUID) (*domain.UserSecurity, error) {
	var security domain.UserSecurity
	if err := r.db.Where("source_user_id = ?", userID).First(&security).Error; err != nil {
		return nil, err
	}
	return &security, nil
}

func (r *userRepository) CreateRefreshToken(token *domain.RefreshToken) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "source_user_id"}, {Name: "device_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"token_hash", "expires_at", "last_used_at"}),
	}).Create(token).Error

}

func (r *userRepository) FindRefreshToken(token string, deviceID string) (*domain.RefreshToken, error) {
	var rt domain.RefreshToken
	err := r.db.Where("token_hash = ? AND device_id = ?", token, deviceID).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *userRepository) RevokeRefreshToken(tokenHash string) error {
	return r.db.Model(&domain.RefreshToken{}).Where("token_hash = ?", tokenHash).Update("revoked_at", time.Now()).Error
}

func (r *userRepository) UpdateUserSecurity(userID uuid.UUID, newPassword string) error {
	return r.db.Model(&domain.UserSecurity{}).Where("source_user_id = ?", userID).Update("password", newPassword).Error
}

func (r *userRepository) AssignRole(userID uuid.UUID, role domain.Role) error {
	return r.db.Create(&domain.ApplicationRole{SourceUserID: userID, Role: role}).Error
}
func (r *userRepository) UpdateRefreshToken(token *domain.RefreshToken) error {
	return r.db.Save(token).Error
}

func (r *userRepository) UpdateUserProfile(profile *domain.UserProfile) error {
	return r.db.Save(profile).Error
}
func (r *userRepository) FindUserProfileByUserID(userID uuid.UUID) (*domain.UserProfile, error) {
	var profile domain.UserProfile
	err := r.db.Preload("Department").Where("source_user_id = ?", userID).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}
