//go:build integration

package config_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/ryu-1905/meow/config"
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

func TestInitDB_Success(t *testing.T) {
	config.ResetDBPoolForTest()
	defer config.ResetDBPoolForTest()

	connStr := getTestDatabaseURL()
	if connStr == "" {
		t.Skip("Bỏ qua test vì không có DATABASE_URL trong môi trường hoặc tệp .env")
	}

	pool, err := config.InitDB(connStr)
	if err != nil {
		t.Fatalf("Không thể khởi tạo cơ sở dữ liệu: %v", err)
	}
	if pool == nil {
		t.Fatal("Pool trả về bị nil dù không có lỗi")
	}

	// Ping thử để kiểm tra pool hoạt động
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		t.Errorf("Ping database thất bại: %v", err)
	}

	// Kiểm tra GetDB()
	if config.GetDB() != pool {
		t.Error("GetDB() không trả về đúng pool vừa khởi tạo")
	}
}
