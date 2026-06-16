package controller

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/ryu-1905/meow/service"
)

// HealthResponse đại diện cho cấu trúc dữ liệu JSON trả về từ API Health Check
type HealthResponse struct {
	Status         string              `json:"status" example:"UP"`
	DatabaseStatus string              `json:"database_status" example:"CONNECTED"`
	ServerInfo     service.ServerInfo  `json:"server_info"`
	AppLogs        []string            `json:"app_logs" example:"[\"2026/05/30 15:00:00 Server starting on port 8080\"]"`
}

// HealthController quản lý các HTTP request liên quan đến sức khỏe hệ thống
type HealthController struct {
	healthService *service.HealthService
}

// NewHealthController khởi tạo một instance mới của HealthController
func NewHealthController(hs *service.HealthService) *HealthController {
	return &HealthController{
		healthService: hs,
	}
}

// GetHealth godoc
// @Summary Kiểm tra trạng thái Server & Logs
// @Description Trả về thông tin chi tiết về tài nguyên hệ thống, bộ nhớ RAM, uptime và 50 dòng log ứng dụng gần nhất.
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (hc *HealthController) GetHealth(c *gin.Context) {
	serverInfo := hc.healthService.GetServerInfo()

	var (
		logs     []string
		err      error
		dbStatus string
		wg       sync.WaitGroup
	)

	wg.Add(2)

	// Lấy log ứng dụng bất đồng bộ
	go func() {
		defer wg.Done()
		logs, err = hc.healthService.GetAppLogs(50)
	}()

	// Kiểm tra trạng thái database bất đồng bộ
	ctx := c.Request.Context()
	go func() {
		defer wg.Done()
		dbStatus = hc.healthService.CheckDatabase(ctx)
	}()

	// Chờ cả hai tác vụ hoàn tất
	wg.Wait()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải log ứng dụng: " + err.Error()})
		return
	}

	response := HealthResponse{
		Status:         "UP",
		DatabaseStatus: dbStatus,
		ServerInfo:     serverInfo,
		AppLogs:        logs,
	}

	c.JSON(http.StatusOK, response)
}
