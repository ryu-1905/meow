//go:build integration

package repository

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

func TestUserRepository_Integration(t *testing.T) {
	connStr := getTestDatabaseURL()
	if connStr == "" {
		t.Skip("Bỏ qua test vì không có DATABASE_URL trong môi trường hoặc tệp .env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("Không thể kết nối database: %v", err)
	}
	defer pool.Close()

	t.Run("Create and Get User Integration Flow", func(t *testing.T) {
		repo := NewUserRepository(pool)

		// Tạo email ngẫu nhiên để tránh trùng lặp dữ liệu
		uniqueEmail := fmt.Sprintf("integration_test_%d@meow.com", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100000))
		expectedName := "Integration User"
		expectedPassword := "hashed_integration_password"

		var userID int64

		// Dọn dẹp dữ liệu sau khi chạy test để tránh làm bẩn database
		defer func() {
			if userID > 0 {
				err := repo.DeleteUser(ctx, userID)
				if err != nil {
					t.Errorf("Lỗi dọn dẹp dữ liệu test: %v", err)
				}
			}
		}()

		// 1. Kiểm tra GetHashPasswordFromEmail với email chưa tồn tại
		_, _, err = repo.GetHashPasswordFromEmail(ctx, uniqueEmail)
		if err == nil {
			t.Fatal("Mong đợi lỗi khi tìm email chưa tồn tại, nhưng không nhận được lỗi")
		}

		// 2. Tạo user mới
		userID, err = repo.CreateUser(ctx, expectedName, uniqueEmail, expectedPassword)
		if err != nil {
			t.Fatalf("Tạo user thất bại: %v", err)
		}
		if userID <= 0 {
			t.Errorf("Mong đợi userID > 0, thực tế nhận được: %d", userID)
		}

		// 3. Lấy lại thông tin user vừa tạo để so khớp
		fetchedID, fetchedHash, err := repo.GetHashPasswordFromEmail(ctx, uniqueEmail)
		if err != nil {
			t.Fatalf("Lấy thông tin user thất bại: %v", err)
		}
		if fetchedID != userID {
			t.Errorf("ID không khớp: mong đợi %d, thực tế %d", userID, fetchedID)
		}
		if fetchedHash != expectedPassword {
			t.Errorf("Mật khẩu băm không khớp: mong đợi %s, thực tế %s", expectedPassword, fetchedHash)
		}
	})
}
