package service

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims định nghĩa các claims chứa thông tin người dùng trong JWT.
type UserClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// JWTService quản lý việc sinh và xác thực JSON Web Token.
type JWTService struct {
	secretKey []byte
}

// NewJWTService khởi tạo một JWTService mới với secret key.
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey: []byte(secretKey),
	}
}

// GenerateAccessToken sinh một Access Token dựa trên User ID.
func (s *JWTService) GenerateAccessToken(userID int64) (string, error) {
	if len(s.secretKey) == 0 {
		return "", errors.New("JWT secret key is empty")
	}

	claims := &UserClaims{
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken giải mã và xác thực tính hợp lệ của token nhận vào.
func (s *JWTService) ValidateToken(tokenStr string) (*UserClaims, error) {
	if len(s.secretKey) == 0 {
		return nil, errors.New("JWT secret key is empty")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (any, error) {
		// Đảm bảo thuật toán ký là HMAC HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
