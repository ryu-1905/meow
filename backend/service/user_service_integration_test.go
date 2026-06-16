//go:build integration

package service_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
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

func TestUserService_Integration(t *testing.T) {
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

	// Khởi tạo các dependencies thực tế
	userRepo := repository.NewUserRepository(pool)
	jwtSecret := "integration-test-secret-key-1234567890"
	jwtService := service.NewJWTService(jwtSecret)
	userService := service.NewUserService(userRepo, jwtService)

	t.Run("Register Success and Validate Token", func(t *testing.T) {
		uniqueEmail := fmt.Sprintf("service_reg_%d@meow.com", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100000))
		name := "Service Integration User"
		password := "password123"

		var userID int64

		// Dọn dẹp dữ liệu kiểm thử
		defer func() {
			if userID > 0 {
				_ = userService.DeleteUser(ctx, userID)
			}
		}()

		// 1. Đăng ký thông qua UserService
		token, err := userService.Register(ctx, name, uniqueEmail, password)
		if err != nil {
			t.Fatalf("Đăng ký qua UserService thất bại: %v", err)
		}
		if token == "" {
			t.Fatal("JWT token trả về rỗng")
		}

		// 2. Xác thực JWT token
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			t.Fatalf("Token trả về không hợp lệ: %v", err)
		}
		if claims.UserID <= 0 {
			t.Errorf("Mong đợi UserID > 0 trong JWT claims, nhận được: %d", claims.UserID)
		}
		userID = claims.UserID

		// 3. Truy vấn database kiểm tra xem user đã tồn tại chưa
		var dbID int64
		var dbName string
		err = pool.QueryRow(ctx, "SELECT id, name FROM users WHERE email = $1", uniqueEmail).Scan(&dbID, &dbName)
		if err != nil {
			t.Fatalf("Không tìm thấy user trong database: %v", err)
		}

		if dbID != claims.UserID {
			t.Errorf("ID trong database (%d) không khớp với ID trong token claims (%d)", dbID, claims.UserID)
		}
		if dbName != name {
			t.Errorf("Tên không khớp: mong đợi %s, thực tế %s", name, dbName)
		}
	})

	t.Run("Register Duplicate Email Error", func(t *testing.T) {
		uniqueEmail := fmt.Sprintf("service_dup_%d@meow.com", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100000))
		name := "Duplicate User"
		password := "password123"

		var userID int64
		defer func() {
			if userID > 0 {
				_ = userService.DeleteUser(ctx, userID)
			}
		}()

		// Đăng ký lần 1
		token, err := userService.Register(ctx, name, uniqueEmail, password)
		if err != nil {
			t.Fatalf("Đăng ký lần 1 thất bại: %v", err)
		}
		claims, err := jwtService.ValidateToken(token)
		if err == nil {
			userID = claims.UserID
		}

		// Đăng ký lần 2 trùng email
		_, err = userService.Register(ctx, "Another Name", uniqueEmail, "anotherpassword")
		if err == nil {
			t.Fatal("Mong đợi đăng ký trùng email sẽ thất bại, nhưng không xảy ra lỗi")
		}
	})

	t.Run("Login Flow", func(t *testing.T) {
		uniqueEmail := fmt.Sprintf("service_login_%d@meow.com", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100000))
		name := "Login Test User"
		password := "mypassword"

		var userID int64
		defer func() {
			if userID > 0 {
				_ = userService.DeleteUser(ctx, userID)
			}
		}()

		// 1. Tạo user bằng cách đăng ký
		token, err := userService.Register(ctx, name, uniqueEmail, password)
		if err != nil {
			t.Fatalf("Đăng ký trước khi login thất bại: %v", err)
		}
		claims, err := jwtService.ValidateToken(token)
		if err == nil {
			userID = claims.UserID
		}

		// 2. Đăng nhập thành công
		token, err = userService.Login(ctx, uniqueEmail, password)
		if err != nil {
			t.Fatalf("Đăng nhập thất bại với thông tin đúng: %v", err)
		}
		if token == "" {
			t.Fatal("Token đăng nhập trả về rỗng")
		}

		// Xác thực token
		claims, err = jwtService.ValidateToken(token)
		if err != nil {
			t.Fatalf("Token đăng nhập không hợp lệ: %v", err)
		}

		var dbID int64
		err = pool.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", uniqueEmail).Scan(&dbID)
		if err != nil {
			t.Fatalf("Không thể tìm thấy user: %v", err)
		}

		if claims.UserID != dbID {
			t.Errorf("UserID trong token (%d) không khớp với DB (%d)", claims.UserID, dbID)
		}

		// 3. Đăng nhập thất bại - sai mật khẩu
		_, err = userService.Login(ctx, uniqueEmail, "wrongpassword")
		if err == nil {
			t.Fatal("Mong đợi đăng nhập thất bại với mật khẩu sai, nhưng không có lỗi")
		}
		if err.Error() != "invalid email or password" {
			t.Errorf("Mong đợi lỗi 'invalid email or password', nhận được: %v", err)
		}

		// 4. Đăng nhập thất bại - email không tồn tại
		_, err = userService.Login(ctx, "nonexistent@meow.com", password)
		if err == nil {
			t.Fatal("Mong đợi đăng nhập thất bại với email không tồn tại, nhưng không có lỗi")
		}
	})
}
