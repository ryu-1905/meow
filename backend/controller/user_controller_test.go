package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ryu-1905/meow/controller"
	"github.com/ryu-1905/meow/service"
	"golang.org/x/crypto/bcrypt"
)

// mockUserRepository giả lập các hành vi của UserRepository.
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

func TestUserController_Register(t *testing.T) {
	// Thiết lập Gin sang Test Mode
	gin.SetMode(gin.TestMode)

	jwtService := service.NewJWTService("test-secret-key")

	t.Run("Register success", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			createUserFunc: func(ctx context.Context, name, email, hashPassword string) (int64, error) {
				return 123, nil
			},
		}

		userService := service.NewUserService(mockRepo, jwtService)
		userController := controller.NewUserController(userService)

		r := gin.Default()
		r.POST("/api/auth/register", userController.Register)

		reqBody, _ := json.Marshal(controller.RegisterRequest{
			Name:     "Test User",
			Email:    "test@meow.com",
			Password: "password123",
		})

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Mong đợi status 200, nhận được: %d", w.Code)
		}

		var resp controller.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("Không thể giải mã response body: %v", err)
		}

		if resp.Token == "" {
			t.Error("Mong đợi token trả về không rỗng")
		}
	})

	t.Run("Register validation error", func(t *testing.T) {
		userService := service.NewUserService(&mockUserRepository{}, jwtService)
		userController := controller.NewUserController(userService)

		r := gin.Default()
		r.POST("/api/auth/register", userController.Register)

		// Gửi email sai định dạng, thiếu name và password
		reqBody, _ := json.Marshal(controller.RegisterRequest{
			Email: "invalid-email",
		})

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Mong đợi status 400 do validation error, nhận được: %d", w.Code)
		}

		if !strings.Contains(w.Body.String(), "không hợp lệ") {
			t.Errorf("Mong đợi thông báo lỗi validation, nhận được: %s", w.Body.String())
		}
	})
}

func TestUserController_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService := service.NewJWTService("test-secret-key")

	// Sinh động bcrypt hash cho "password123" để đảm bảo so khớp mật khẩu chính xác
	hashedPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedPassword := string(hashedPasswordBytes)

	t.Run("Login success", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			getHashPasswordFromEmailFunc: func(ctx context.Context, email string) (int64, string, error) {
				return 456, hashedPassword, nil
			},
		}

		userService := service.NewUserService(mockRepo, jwtService)
		userController := controller.NewUserController(userService)

		r := gin.Default()
		r.POST("/api/auth/login", userController.Login)

		reqBody, _ := json.Marshal(controller.LoginRequest{
			Email:    "test@meow.com",
			Password: "password123",
		})

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Mong đợi status 200, nhận được: %d. Body: %s", w.Code, w.Body.String())
		}

		var resp controller.AuthResponse
		_ = json.Unmarshal(w.Body.Bytes(), &resp)

		if resp.Token == "" {
			t.Error("Mong đợi token trả về không rỗng")
		}
	})

	t.Run("Login unauthorized - wrong password", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			getHashPasswordFromEmailFunc: func(ctx context.Context, email string) (int64, string, error) {
				return 456, hashedPassword, nil
			},
		}

		userService := service.NewUserService(mockRepo, jwtService)
		userController := controller.NewUserController(userService)

		r := gin.Default()
		r.POST("/api/auth/login", userController.Login)

		reqBody, _ := json.Marshal(controller.LoginRequest{
			Email:    "test@meow.com",
			Password: "wrongpassword",
		})

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Mong đợi status 401 do sai mật khẩu, nhận được: %d", w.Code)
		}

		if !strings.Contains(w.Body.String(), "Sai địa chỉ email hoặc mật khẩu") {
			t.Errorf("Mong đợi thông báo lỗi sai thông tin đăng nhập, nhận được: %s", w.Body.String())
		}
	})
}
