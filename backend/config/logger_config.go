package config

import (
	"io"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

// InitLogger cấu hình ghi log kép (terminal + file app.log) và thiết lập slog default logger.
// Trả về log file để main có thể defer close nếu cần.
func InitLogger() *os.File {
	var logWriter io.Writer = os.Stdout
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Không thể khởi tạo file log", slog.Any("error", err))
	} else {
		// Thiết lập gin writer và log writer để ghi ra cả terminal và file log
		gin.DefaultWriter = io.MultiWriter(os.Stdout, logFile)
		logWriter = io.MultiWriter(os.Stdout, logFile)
	}

	// Thiết lập slog default logger
	logger := slog.New(slog.NewTextHandler(logWriter, nil))
	slog.SetDefault(logger)

	return logFile
}
