package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ryu-1905/meow/config"
	"github.com/ryu-1905/meow/controller"
	"github.com/ryu-1905/meow/service"
)

// TestGetHealth_Disconnected kiểm tra API Health Check khi Database bị ngắt kết nối
func TestGetHealth_Disconnected(t *testing.T) {
	// Chuyển Gin sang chế độ test
	gin.SetMode(gin.TestMode)

	// Khởi tạo service không có db connection (nil)
	hs := service.NewHealthService(nil)
	hc := controller.NewHealthController(hs)

	// Thiết lập router test
	r := gin.New()
	r.GET("/api/health", hc.GetHealth)

	// Tạo request giả lập
	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatalf("Không thể tạo request: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Mong đợi status 200, nhưng nhận được %d", w.Code)
	}

	var response controller.HealthResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Không thể parse JSON phản hồi: %v", err)
	}

	if response.Status != "UP" {
		t.Errorf("Mong đợi status 'UP', nhưng nhận được %s", response.Status)
	}

	if response.DatabaseStatus != "DISCONNECTED" {
		t.Errorf("Mong đợi database_status 'DISCONNECTED', nhưng nhận được %s", response.DatabaseStatus)
	}
}

// TestGetHealth_Connected kiểm tra API Health Check khi kết nối thành công tới Database thật
func TestGetHealth_Connected(t *testing.T) {
	// Chuyển Gin sang chế độ test
	gin.SetMode(gin.TestMode)

	// Thử load file .env từ thư mục cha hoặc thư mục hiện tại để lấy DATABASE_URL
	_ = godotenv.Load("../.env", ".env")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("Bỏ qua TestGetHealth_Connected vì DATABASE_URL không được thiết lập trong môi trường hoặc file .env")
	}

	// Khởi tạo Database Pool thật
	dbPool, err := config.InitDB(dbURL)
	if err != nil {
		t.Fatalf("Không thể kết nối tới cơ sở dữ liệu thật: %v", err)
	}
	defer dbPool.Close()

	hs := service.NewHealthService(dbPool)
	hc := controller.NewHealthController(hs)

	// Thiết lập router test
	r := gin.New()
	r.GET("/api/health", hc.GetHealth)

	// Tạo request giả lập
	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatalf("Không thể tạo request: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Mong đợi status 200, nhưng nhận được %d", w.Code)
	}

	var response controller.HealthResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Không thể parse JSON phản hồi: %v", err)
	}

	if response.Status != "UP" {
		t.Errorf("Mong đợi status 'UP', nhưng nhận được %s", response.Status)
	}

	if response.DatabaseStatus != "CONNECTED" {
		t.Errorf("Mong đợi database_status 'CONNECTED', nhưng nhận được %s", response.DatabaseStatus)
	}
}
