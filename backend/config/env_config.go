package config

import (
	"log/slog"

	"github.com/joho/godotenv"
)

// LoadEnv tải các biến môi trường từ tệp .env.
// Nếu không tìm thấy tệp .env, hàm sẽ log info và tiếp tục sử dụng biến môi trường mặc định.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found, using default environment variables")
	}
}
