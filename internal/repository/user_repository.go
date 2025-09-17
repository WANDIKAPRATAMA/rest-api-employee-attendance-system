package repository

import (
	"employee-attendance-system/internal/entity/domain"
	"employee-attendance-system/internal/entity/dto"
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
	IsUserExist(userID uuid.UUID) (bool, error)
	FindAllUsers(req dto.ListUsersRequest) ([]*domain.UserProfile, int64, error)

	CountEmployeesPerDepartment() (map[string]int, error)
	CountTodayRegistrations(today time.Time) (int, error)
}

type userRepository struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewUserRepository(db *gorm.DB, log *logrus.Logger) UserRepository {
	return &userRepository{db: db, log: log}
}

func (r *userRepository) CountEmployeesPerDepartment() (map[string]int, error) {
	var results []struct {
		DeptName string
		Count    int
	}
	err := r.db.Model(&domain.UserProfile{}).
		Joins("LEFT JOIN departments ON departments.id = user_profiles.department_id").
		Where("user_profiles.deleted_at IS NULL AND departments.deleted_at IS NULL").
		Select("departments.department_name AS dept_name, COUNT(user_profiles.id) AS count").
		Group("departments.department_name").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	employeesPerDept := make(map[string]int)
	for _, r := range results {
		employeesPerDept[r.DeptName] = r.Count
	}
	return employeesPerDept, nil
}

func (r *userRepository) CountTodayRegistrations(today time.Time) (int, error) {
	var count int64
	err := r.db.Model(&domain.UserProfile{}).
		Where("created_at >= ? AND created_at < ? AND deleted_at IS NULL", today, today.Add(24*time.Hour)).
		Count(&count).Error
	return int(count), err
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
func (r *userRepository) IsUserExist(userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.
		Model(&domain.User{}).
		Select("1").
		Where("id = ?", userID).
		Limit(1).
		Scan(&exists).Error

	if err != nil {
		return false, err
	}
	return exists, nil
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
	err := r.db.
		Model(&domain.UserProfile{}).
		Preload("Department").
		Preload("ApplicationRole").
		Where("source_user_id = ?", userID).
		First(&profile).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *userRepository) FindAllUsers(req dto.ListUsersRequest) ([]*domain.UserProfile, int64, error) {
	var users []*domain.UserProfile

	query := r.db.Model(&domain.UserProfile{}).
		Preload("Department").
		Preload("ApplicationRole").
		Where("user_profiles.deleted_at IS NULL")

	// Dynamic filters
	if req.Email != "" {
		query = query.Joins("JOIN users u ON u.id = user_profiles.source_user_id").
			Where("u.email LIKE ?", "%"+req.Email+"%")
	}
	if req.Status != "" {
		query = query.Joins("JOIN users u ON u.id = user_profiles.source_user_id").
			Where("u.status = ?", req.Status)
	}
	if req.DepartmentID != nil {
		query = query.Where("department_id = ?", *req.DepartmentID)
	}
	if req.CreatedAtStart != nil {
		query = query.Where("user_profiles.created_at >= ?", req.CreatedAtStart)
	}
	if req.CreatedAtEnd != nil {
		query = query.Where("user_profiles.created_at <= ?", req.CreatedAtEnd)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
