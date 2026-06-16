package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ryu-1905/meow/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService quản lý các logic nghiệp vụ liên quan đến người dùng (đăng ký, đăng nhập).
type UserService struct {
	userRepo   repository.UserRepository
	jwtService *JWTService
}

// NewUserService khởi tạo một UserService mới.
func NewUserService(userRepo repository.UserRepository, jwtService *JWTService) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Register thực hiện băm mật khẩu người dùng, tạo bản ghi mới trong cơ sở dữ liệu và sinh Access Token.
func (s *UserService) Register(ctx context.Context, name, email, password string) (string, error) {
	if name == "" || email == "" || password == "" {
		return "", errors.New("name, email and password cannot be empty")
	}

	// Băm mật khẩu sử dụng bcrypt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Tạo người dùng trong database
	userID, err := s.userRepo.CreateUser(ctx, name, email, string(hashedBytes))
	if err != nil {
		return "", err
	}

	// Sinh access token sau khi đăng ký thành công
	token, err := s.jwtService.GenerateAccessToken(userID)
	if err != nil {
		return "", fmt.Errorf("user created but failed to generate access token: %w", err)
	}

	return token, nil
}

// Login kiểm tra thông tin email, so khớp mật khẩu và sinh ra access token nếu hợp lệ.
func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("email and password cannot be empty")
	}

	// Lấy ID và mật khẩu đã băm từ email
	userID, hashPassword, err := s.userRepo.GetHashPasswordFromEmail(ctx, email)
	if err != nil {
		return "", err
	}

	// So sánh mật khẩu thô gửi lên với mật khẩu đã băm
	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Sinh access token từ User ID
	token, err := s.jwtService.GenerateAccessToken(userID)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return token, nil
}

// DeleteUser xóa người dùng trong hệ thống dựa trên ID.
func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid user id")
	}
	return s.userRepo.DeleteUser(ctx, id)
}
