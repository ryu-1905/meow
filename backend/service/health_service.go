package service

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ServerInfo chứa thông tin cấu hình và tài nguyên của server
type ServerInfo struct {
	OS           string    `json:"os" example:"windows"`
	Architecture string    `json:"architecture" example:"amd64"`
	GoVersion    string    `json:"go_version" example:"go1.26.3"`
	NumCPU       int       `json:"num_cpu" example:"8"`
	NumGoroutine int       `json:"num_goroutine" example:"4"`
	Uptime       string    `json:"uptime" example:"1h2m3s"`
	MemoryUsage  MemStatus `json:"memory_usage"`
}

// MemStatus chứa thông tin về bộ nhớ RAM đang được ứng dụng sử dụng
type MemStatus struct {
	Alloc      string `json:"alloc" example:"2.50 MB"`
	TotalAlloc string `json:"total_alloc" example:"10.20 MB"`
	Sys        string `json:"sys" example:"15.30 MB"`
	NumGC      uint32 `json:"num_gc" example:"1"`
}

var startTime time.Time

func init() {
	startTime = time.Now()
}

// HealthService quản lý logic nghiệp vụ kiểm tra trạng thái hệ thống
type HealthService struct {
	dbPool *pgxpool.Pool
}

// NewHealthService khởi tạo một instance mới của HealthService
func NewHealthService(dbPool *pgxpool.Pool) *HealthService {
	return &HealthService{
		dbPool: dbPool,
	}
}

// CheckDatabase kiểm tra xem kết nối tới database có hoạt động hay không
func (s *HealthService) CheckDatabase(ctx context.Context) string {
	if s.dbPool == nil {
		return "DISCONNECTED"
	}
	if err := s.dbPool.Ping(ctx); err != nil {
		return "DISCONNECTED (Ping error: " + err.Error() + ")"
	}
	return "CONNECTED"
}

// GetServerInfo thu thập thông số chi tiết của hệ thống và runtime Go
func (s *HealthService) GetServerInfo() ServerInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return ServerInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		GoVersion:    runtime.Version(),
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		Uptime:       time.Since(startTime).Truncate(time.Second).String(),
		MemoryUsage: MemStatus{
			Alloc:      fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024),
			TotalAlloc: fmt.Sprintf("%.2f MB", float64(m.TotalAlloc)/1024/1024),
			Sys:        fmt.Sprintf("%.2f MB", float64(m.Sys)/1024/1024),
			NumGC:      m.NumGC,
		},
	}
}

// GetAppLogs đọc N dòng log cuối cùng từ file log của ứng dụng (app.log)
func (s *HealthService) GetAppLogs(linesCount int) ([]string, error) {
	file, err := os.Open("app.log")
	if err != nil {
		if os.IsNotExist(err) {
			return []string{"Không tìm thấy log file. Đảm bảo app.log đã được tạo và ghi log hoạt động."}, nil
		}
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Lấy ra N dòng log cuối cùng
	start := 0
	if len(lines) > linesCount {
		start = len(lines) - linesCount
	}

	return lines[start:], nil
}
