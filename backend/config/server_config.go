package config

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

// StartServer khởi chạy máy chủ Gin trên PORT tương ứng được định nghĩa trong biến môi trường.
// Nếu không cấu hình PORT, mặc định server sẽ chạy ở port 8080.
func StartServer(r *gin.Engine) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Server starting", slog.String("port", port))
	if err := r.Run(":" + port); err != nil {
		slog.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
