package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type JWTConfig struct {
	AccesTokenSecretKey   string
	RefreshTokenSecretKey string
}
type ValidatedMethod string

const (
	AccessToken  ValidatedMethod = "access_token"
	RefreshToken ValidatedMethod = "refresh_token"
)

func NewJWTCfg(viper *viper.Viper) *JWTConfig {
	return &JWTConfig{
		AccesTokenSecretKey:   viper.GetString("jwt.accesTokenSecret"),
		RefreshTokenSecretKey: viper.GetString("jwt.refreshTokenSecret"),
	}
}

// GenerateAccessToken membuat JWT access token
func (j *JWTConfig) GenerateToken(ctx context.Context, userID uuid.UUID, email, role string, method ValidatedMethod) (string, error) {
	var expires time.Duration
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"role":    role,
	}

	if method == RefreshToken {
		expires = 30 * 24 * time.Hour
		claims["exp"] = time.Now().Add(expires).Unix()
	} else {
		expires = 7 * 24 * time.Hour
		claims["exp"] = time.Now().Add(expires).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// fix: ambil 2 value
	var (
		result string
		err    error
	)
	if method == RefreshToken {
		result, err = token.SignedString([]byte(j.RefreshTokenSecretKey))
	} else {
		result, err = token.SignedString([]byte(j.AccesTokenSecretKey))
	}
	if err != nil {
		return "", err
	}

	// TODO: REMOVE Implmentation With Redis simpan token ke redis (pakai string langsung, tanpa pointer)
	// if _, err := j.RedisClient.Set(ctx, result, userID.String(), expires).Result(); err != nil {
	// 	return "", err
	// }

	return result, nil
}

// ValidateToken memverifikasi token dan validasi ke redis
func (j *JWTConfig) ValidateToken(ctx context.Context, tokenString string, method ValidatedMethod) (*jwt.Token, error) {
	// validasi ke redis apakah token exist
	// val, err := j.RedisClient.Get(ctx, tokenString).Result()
	// if err == redis.Nil {
	// 	return nil, fmt.Errorf("token not found in redis")
	// } else if err != nil {
	// 	return nil, err
	// }

	// if val != tokenString {
	// 	return nil, fmt.Errorf("invalid token (not matched in redis)")
	// }

	// parse JWT
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		if method == AccessToken {
			return []byte(j.AccesTokenSecretKey), nil
		}
		return []byte(j.RefreshTokenSecretKey), nil
	})
}
