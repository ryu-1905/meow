package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// mockUserRepository giả lập các hành vi của UserRepository phục vụ cho việc test độc lập.
type mockUserRepository struct {
	createUserFunc               func(ctx context.Context, name, email, hashPassword string) (int64, error)
	getHashPasswordFromEmailFunc func(ctx context.Context, email string) (int64, string, error)
	deleteUserFunc               func(ctx context.Context, id int64) error
}

func (m *mockUserRepository) CreateUser(ctx context.Context, name, email, hashPassword string) (int64, error) {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, name, email, hashPassword)
	}
	return 0, nil
}

func (m *mockUserRepository) GetHashPasswordFromEmail(ctx context.Context, email string) (int64, string, error) {
	if m.getHashPasswordFromEmailFunc != nil {
		return m.getHashPasswordFromEmailFunc(ctx, email)
	}
	return 0, "", nil
}

func (m *mockUserRepository) DeleteUser(ctx context.Context, id int64) error {
	if m.deleteUserFunc != nil {
		return m.deleteUserFunc(ctx, id)
	}
	return nil
}

func TestJWTService_GenerateAndValidateToken(t *testing.T) {
	secret := "my-super-secret-key-12345"
	jwtService := NewJWTService(secret)

	var userID int64 = 42

	// 1. Kiểm tra Generate Token
	tokenStr, err := jwtService.GenerateAccessToken(userID)
	if err != nil {
		t.Fatalf("Không thể tạo token: %v", err)
	}

	if tokenStr == "" {
		t.Fatal("Token được tạo ra bị rỗng")
	}

	// 2. Kiểm tra Validate Token
	claims, err := jwtService.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("Không thể xác thực token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID giải mã không khớp: mong đợi %d, thực tế %d", userID, claims.UserID)
	}
}

func TestJWTService_InvalidSignature(t *testing.T) {
	secret1 := "secret-number-one"
	secret2 := "secret-number-two"

	jwtService1 := NewJWTService(secret1)
	jwtService2 := NewJWTService(secret2)

	tokenStr, err := jwtService1.GenerateAccessToken(99)
	if err != nil {
		t.Fatalf("Không thể tạo token: %v", err)
	}

	// Sử dụng service 2 (có secret key khác) để xác thực token của service 1
	_, err = jwtService2.ValidateToken(tokenStr)
	if err == nil {
		t.Fatal("Mong đợi lỗi xác thực do sai chữ ký bí mật nhưng không có lỗi xảy ra")
	}

	if !strings.Contains(err.Error(), "signature is invalid") && !strings.Contains(err.Error(), "failed to parse") {
		t.Errorf("Mong đợi lỗi liên quan đến signature, thực tế nhận được: %v", err)
	}
}

func TestUserService_Register(t *testing.T) {
	jwtService := NewJWTService("test-secret-key")

	t.Run("Register success", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			createUserFunc: func(ctx context.Context, name, email, hashPassword string) (int64, error) {
				if name != "John" || email != "john@meow.com" {
					return 0, errors.New("unexpected arguments")
				}
				// Kiểm tra mật khẩu đã được băm bằng bcrypt
				err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte("password123"))
				if err != nil {
					return 0, errors.New("password not hashed correctly")
				}
				return 100, nil
			},
		}

		userService := NewUserService(mockRepo, jwtService)
		token, err := userService.Register(context.Background(), "John", "john@meow.com", "password123")
		if err != nil {
			t.Fatalf("Mong đợi không có lỗi, thực tế nhận lỗi: %v", err)
		}

		if token == "" {
			t.Fatal("Access Token trả về bị rỗng")
		}

		// Xác thực token được trả ra chứa UserID = 100
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			t.Fatalf("Token sinh ra không hợp lệ: %v", err)
		}

		if claims.UserID != 100 {
			t.Errorf("Mong đợi UserID trong token là 100, thực tế nhận: %d", claims.UserID)
		}
	})

	t.Run("Validation error - empty arguments", func(t *testing.T) {
		mockRepo := &mockUserRepository{}
		userService := NewUserService(mockRepo, jwtService)

		_, err := userService.Register(context.Background(), "", "john@meow.com", "password123")
		if err == nil || !strings.Contains(err.Error(), "cannot be empty") {
			t.Errorf("Mong đợi lỗi validate tham số trống, nhận được: %v", err)
		}
	})

	t.Run("Database error on CreateUser", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			createUserFunc: func(ctx context.Context, name, email, hashPassword string) (int64, error) {
				return 0, errors.New("db unique constraint error")
			},
		}

		userService := NewUserService(mockRepo, jwtService)
		_, err := userService.Register(context.Background(), "John", "john@meow.com", "password123")
		if err == nil || !strings.Contains(err.Error(), "db unique") {
			t.Errorf("Mong đợi chuyển tiếp lỗi từ database, nhận được: %v", err)
		}
	})
}

func TestUserService_Login(t *testing.T) {
	jwtService := NewJWTService("test-secret-key")
	var mockUserID int64 = 55
	rawPassword := "mypassword"
	hashedPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	mockHashedPassword := string(hashedPasswordBytes)

	t.Run("Login success", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			getHashPasswordFromEmailFunc: func(ctx context.Context, email string) (int64, string, error) {
				if email != "test@meow.com" {
					return 0, "", errors.New("user not found")
				}
				return mockUserID, mockHashedPassword, nil
			},
		}

		userService := NewUserService(mockRepo, jwtService)
		token, err := userService.Login(context.Background(), "test@meow.com", rawPassword)
		if err != nil {
			t.Fatalf("Đăng nhập thất bại: %v", err)
		}

		if token == "" {
			t.Fatal("Token trả ra bị rỗng")
		}

		// Xác thực token
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			t.Fatalf("Token không hợp lệ: %v", err)
		}

		if claims.UserID != mockUserID {
			t.Errorf("UserID trong token không khớp: mong đợi %d, thực tế %d", mockUserID, claims.UserID)
		}
	})

	t.Run("Validation error - empty credentials", func(t *testing.T) {
		mockRepo := &mockUserRepository{}
		userService := NewUserService(mockRepo, jwtService)

		_, err := userService.Login(context.Background(), "", rawPassword)
		if err == nil || !strings.Contains(err.Error(), "cannot be empty") {
			t.Errorf("Mong đợi lỗi validate, nhận được: %v", err)
		}
	})

	t.Run("User not found in database", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			getHashPasswordFromEmailFunc: func(ctx context.Context, email string) (int64, string, error) {
				return 0, "", errors.New("user not found")
			},
		}

		userService := NewUserService(mockRepo, jwtService)
		_, err := userService.Login(context.Background(), "unknown@meow.com", rawPassword)
		if err == nil || !strings.Contains(err.Error(), "user not found") {
			t.Errorf("Mong đợi lỗi user not found, nhận được: %v", err)
		}
	})

	t.Run("Invalid password", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			getHashPasswordFromEmailFunc: func(ctx context.Context, email string) (int64, string, error) {
				return mockUserID, mockHashedPassword, nil
			},
		}

		userService := NewUserService(mockRepo, jwtService)
		_, err := userService.Login(context.Background(), "test@meow.com", "wrongpassword")
		if err == nil || !strings.Contains(err.Error(), "invalid email or password") {
			t.Errorf("Mong đợi lỗi mật khẩu không khớp, nhận được: %v", err)
		}
	})
}
