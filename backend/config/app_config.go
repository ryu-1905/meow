package config

import (
	"log/slog"
	"os"
)

// ConfigApp thực hiện toàn bộ các bước cấu hình, tiêm phụ thuộc, thiết lập router và khởi chạy ứng dụng.
func ConfigApp() {
	// Khởi tạo logger
	logFile := InitLogger()
	if logFile != nil {
		defer logFile.Close()
	}

	// Tải cấu hình từ file .env
	LoadEnv()

	// Khởi tạo cơ sở dữ liệu từ biến môi trường
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		slog.Error("DATABASE_URL environment variable is not set")
		os.Exit(1)
	}

	dbPool, err := InitDB(dbURL)
	if err != nil {
		slog.Error("Không thể khởi tạo cơ sở dữ liệu", slog.Any("error", err))
		os.Exit(1)
	}
	defer dbPool.Close()

	// Khởi tạo Dependency Injection Container
	container := NewContainer(dbPool)

	// Thiết lập router
	r := SetupRouter(container)

	// Khởi chạy máy chủ
	StartServer(r)
}
