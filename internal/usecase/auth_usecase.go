package usecase

import (
	"context"
	"crypto/rand"
	"employee-attendance-system/internal/entity/domain"
	"employee-attendance-system/internal/entity/dto"
	"employee-attendance-system/internal/repository"
	utils "employee-attendance-system/internal/util"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCase interface {
	Signup(ctx context.Context, email, password, fullName string) (*domain.User, error)
	Signin(ctx context.Context, email, password string, deviceID *string) (string, string, *dto.UserResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	RefreshToken(ctx context.Context, refreshToken string, deviceID string) (string, string, error) // newAccessToken
	ChangeRole(ctx context.Context, userID uuid.UUID, role string) error
	Signout(ctx context.Context, tokenHash string) error
}

type authUseCase struct {
	repo     repository.UserRepository
	validate *validator.Validate
	log      *logrus.Logger
	config   *viper.Viper
	jwtUtils *utils.JWTConfig
}

func NewAuthUseCase(
	repo repository.UserRepository,
	log *logrus.Logger,
	validate *validator.Validate,
	config *viper.Viper,
	jwtUtils *utils.JWTConfig,
) AuthUseCase {
	return &authUseCase{repo: repo, log: log, validate: validate, config: config,
		jwtUtils: jwtUtils}

}
func GenerateEmployeeCode() string {
	return fmt.Sprintf("EMP-%d", time.Now().UnixNano())
}

func (u *authUseCase) Signup(ctx context.Context, email, password, fullName string) (*domain.User, error) {

	exist, err := u.repo.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	code := GenerateEmployeeCode()
	user := &domain.User{Email: email, Status: "active", EmailVerified: true}
	profile := &domain.UserProfile{FullName: fullName,
		EmployeeCode: code}
	security := &domain.UserSecurity{Password: string(hashedPassword)}
	role := &domain.ApplicationRole{Role: "user"}

	if err := u.repo.CreateUser(user, profile, security, role); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *authUseCase) Signin(ctx context.Context, email, password string, deviceID *string) (string, string, *dto.UserResponse, error) {
	user, err := u.repo.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", nil, fmt.Errorf("invalid email or password")
		}
		return "", "", nil, err
	}
	if user == nil {
		return "", "", nil, fmt.Errorf("invalid email or password")
	}

	security, err := u.repo.FindUserSecurityByUserID(user.ID)
	if err != nil {
		return "", "", nil, fmt.Errorf("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(security.Password), []byte(password)); err != nil {
		return "", "", nil, fmt.Errorf("invalid email or password")
	}

	role, err := u.repo.FindUserRoleByUserID(user.ID)
	if err != nil {
		return "", "", nil, err
	}

	accessToken, err := u.jwtUtils.GenerateToken(ctx, user.ID, user.Email, string(role), utils.AccessToken)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err := u.jwtUtils.GenerateToken(ctx, user.ID, user.Email, string(role), utils.RefreshToken)
	if err != nil {
		return "", "", nil, err
	}

	refresh := &domain.RefreshToken{
		SourceUserID: user.ID,
		TokenHash:    refreshToken,
		ExpiresAt:    time.Now().Add(48 * time.Hour), // hati-hati jangan * 24 dua kali
	}
	if deviceID != nil {
		refresh.DeviceID = *deviceID
	}

	if err := u.repo.CreateRefreshToken(refresh); err != nil {
		return "", "", nil, err
	}

	profile, err := u.repo.FindUserProfileByUserID(user.ID)
	if err != nil {
		return "", "", nil, err
	}
	r := mapToUserResponse(profile)
	r.Email = user.Email

	return accessToken, refreshToken, r, nil
}

func (u *authUseCase) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	security, err := u.repo.FindUserSecurityByUserID(userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(security.Password), []byte(oldPassword)); err != nil {
		return err
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return u.repo.UpdateUserSecurity(userID, string(hashedNewPassword))
}

func (u *authUseCase) RefreshToken(ctx context.Context, refreshToken string, deviceID string) (string, string, error) {

	u.log.Println("device id", deviceID, "refresh token", refreshToken)
	storedToken, err := u.repo.FindRefreshToken(refreshToken, deviceID)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	u.log.Println("stored token", storedToken)
	if storedToken.RevokedAt != nil && !storedToken.RevokedAt.IsZero() {
		return "", "", fmt.Errorf("refresh token revoked")
	}

	if time.Now().After(storedToken.ExpiresAt) {
		return "", "", fmt.Errorf("refresh token expired")
	}

	// Ambil user
	user, err := u.repo.FindUserByID(storedToken.SourceUserID)
	if err != nil {
		return "", "", fmt.Errorf("user not found")
	}

	role, err := u.repo.FindUserRoleByUserID(user.ID)
	if err != nil {
		return "", "", err
	}

	// Generate access token baru
	accessToken, err := u.jwtUtils.GenerateToken(ctx, user.ID, user.Email, string(role), utils.AccessToken)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token baru (opsional, best practice rotate)
	refreshTokenBytes := make([]byte, 32)
	rand.Read(refreshTokenBytes)
	newRefreshToken := hex.EncodeToString(refreshTokenBytes)

	storedToken.TokenHash = newRefreshToken
	storedToken.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
	storedToken.LastUsedAt = time.Now()

	if err := u.repo.UpdateRefreshToken(storedToken); err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (u *authUseCase) ChangeRole(ctx context.Context, userID uuid.UUID, role string) error {
	r := domain.Role(role)
	return u.repo.AssignRole(userID, r)
}

func (u *authUseCase) Signout(ctx context.Context, tokenHash string) error {
	return u.repo.RevokeRefreshToken(tokenHash)
}
