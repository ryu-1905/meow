//go:build integration

package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/ryu-1905/meow/controller"
	"github.com/ryu-1905/meow/repository"
	"github.com/ryu-1905/meow/service"
)

func getTestDatabaseURL() string {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		return dbURL
	}

	dir, err := os.Getwd()
	if err == nil {
		for {
			envPath := filepath.Join(dir, ".env")
			if _, err := os.Stat(envPath); err == nil {
				_ = godotenv.Load(envPath)
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return os.Getenv("DATABASE_URL")
}

func TestUserController_Integration(t *testing.T) {
	connStr := getTestDatabaseURL()
	if connStr == "" {
		t.Skip("Bỏ qua test vì không có DATABASE_URL trong môi trường hoặc tệp .env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("Không thể kết nối database: %v", err)
	}
	defer pool.Close()

	// Khởi tạo các dependencies thật
	userRepo := repository.NewUserRepository(pool)
	jwtSecret := "controller-integration-jwt-secret-key-987654321"
	jwtService := service.NewJWTService(jwtSecret)
	userService := service.NewUserService(userRepo, jwtService)
	userController := controller.NewUserController(userService)

	// Thiết lập Gin Test Mode
	gin.SetMode(gin.TestMode)

	t.Run("HTTP Register Success", func(t *testing.T) {
		uniqueEmail := fmt.Sprintf("ctrl_reg_%d@meow.com", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100000))
		name := "Controller Integration User"
		password := "securepassword123"

		var userID int64
		defer func() {
			if userID > 0 {
				_ = userService.DeleteUser(ctx, userID)
			}
		}()

		r := gin.New()
		r.POST("/api/auth/register", userController.Register)

		reqBody, _ := json.Marshal(controller.RegisterRequest{
			Name:     name,
			Email:    uniqueEmail,
			Password: password,
		})

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Mong đợi status 200, nhận được: %d. Body: %s", w.Code, w.Body.String())
		}

		var resp controller.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("Không thể giải mã response body: %v", err)
		}

		if resp.Token == "" {
			t.Fatal("JWT Token trả về rỗng")
		}

		claims, err := jwtService.ValidateToken(resp.Token)
		if err == nil {
			userID = claims.UserID
		}

		// Xác thực dữ liệu đã chèn trong DB
		var dbName string
		err = pool.QueryRow(ctx, "SELECT name FROM users WHERE email = $1", uniqueEmail).Scan(&dbName)
		if err != nil {
			t.Fatalf("Không tìm thấy user trong database: %v", err)
		}
		if dbName != name {
			t.Errorf("Tên trong DB không khớp: mong đợi %s, thực tế %s", name, dbName)
		}
	})

	t.Run("HTTP Login Flow", func(t *testing.T) {
		uniqueEmail := fmt.Sprintf("ctrl_login_%d@meow.com", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100000))
		name := "Controller Login User"
		password := "password321"

		var userID int64
		defer func() {
			if userID > 0 {
				_ = userService.DeleteUser(ctx, userID)
			}
		}()

		// 1. Tạo trước user thông qua service
		token, err := userService.Register(ctx, name, uniqueEmail, password)
		if err != nil {
			t.Fatalf("Tạo user phục vụ kiểm thử login thất bại: %v", err)
		}
		claims, err := jwtService.ValidateToken(token)
		if err == nil {
			userID = claims.UserID
		}

		r := gin.New()
		r.POST("/api/auth/login", userController.Login)

		// 2. Đăng nhập thành công
		reqBody, _ := json.Marshal(controller.LoginRequest{
			Email:    uniqueEmail,
			Password: password,
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Mong đợi status 200, nhận được: %d. Body: %s", w.Code, w.Body.String())
		}

		var resp controller.AuthResponse
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		if resp.Token == "" {
			t.Fatal("Token trả về khi login bị rỗng")
		}

		// 3. Đăng nhập thất bại - sai mật khẩu
		wrongReqBody, _ := json.Marshal(controller.LoginRequest{
			Email:    uniqueEmail,
			Password: "wrongpassword",
		})
		reqWrong, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(wrongReqBody))
		reqWrong.Header.Set("Content-Type", "application/json")

		wWrong := httptest.NewRecorder()
		r.ServeHTTP(wWrong, reqWrong)

		if wWrong.Code != http.StatusUnauthorized {
			t.Errorf("Mong đợi status 401, nhận được: %d. Body: %s", wWrong.Code, wWrong.Body.String())
		}
		if !strings.Contains(wWrong.Body.String(), "Sai địa chỉ email hoặc mật khẩu") {
			t.Errorf("Mong đợi thông điệp lỗi chính xác, nhận được: %s", wWrong.Body.String())
		}
	})
}
