package service_test

import (
	"context"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/ryu-1905/meow/config"
	"github.com/ryu-1905/meow/service"
)

func TestHealthService_CheckDatabase(t *testing.T) {
	// Case 1: dbPool is nil
	t.Run("dbPool is nil", func(t *testing.T) {
		s := service.NewHealthService(nil)
		status := s.CheckDatabase(context.Background())
		if status != "DISCONNECTED" {
			t.Errorf("Mong đợi status 'DISCONNECTED', nhưng nhận được %q", status)
		}
	})

	// Case 2: dbPool is connected (if DATABASE_URL is available)
	t.Run("dbPool is connected", func(t *testing.T) {
		// Thử load file .env từ thư mục cha hoặc thư mục hiện tại
		_ = godotenv.Load("../.env", ".env")

		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			t.Skip("Bỏ qua test vì DATABASE_URL không được thiết lập trong môi trường hoặc file .env")
		}

		dbPool, err := config.InitDB(dbURL)
		if err != nil {
			t.Fatalf("Không thể khởi tạo DB Pool: %v", err)
		}

		s := service.NewHealthService(dbPool)
		status := s.CheckDatabase(context.Background())
		if status != "CONNECTED" {
			t.Errorf("Mong đợi status 'CONNECTED', nhưng nhận được %q", status)
		}
	})
}

func TestHealthService_GetServerInfo(t *testing.T) {
	s := service.NewHealthService(nil)
	info := s.GetServerInfo()

	if info.OS != runtime.GOOS {
		t.Errorf("Mong đợi OS %q, nhưng nhận được %q", runtime.GOOS, info.OS)
	}

	if info.Architecture != runtime.GOARCH {
		t.Errorf("Mong đợi Architecture %q, nhưng nhận được %q", runtime.GOARCH, info.Architecture)
	}

	if info.GoVersion != runtime.Version() {
		t.Errorf("Mong đợi GoVersion %q, nhưng nhận được %q", runtime.Version(), info.GoVersion)
	}

	if info.NumCPU != runtime.NumCPU() {
		t.Errorf("Mong đợi NumCPU %d, nhưng nhận được %d", runtime.NumCPU(), info.NumCPU)
	}

	if info.NumGoroutine <= 0 {
		t.Errorf("Mong đợi NumGoroutine > 0, nhưng nhận được %d", info.NumGoroutine)
	}

	if info.Uptime == "" {
		t.Error("Mong đợi Uptime không rỗng")
	}

	if info.MemoryUsage.Alloc == "" || info.MemoryUsage.TotalAlloc == "" || info.MemoryUsage.Sys == "" {
		t.Error("Mong đợi thông tin RAM (Alloc, TotalAlloc, Sys) không được rỗng")
	}
}

func TestHealthService_GetAppLogs(t *testing.T) {
	// Backup file app.log hiện tại nếu có
	hasBackup := false
	if _, err := os.Stat("app.log"); err == nil {
		err := os.Rename("app.log", "app.log.bak")
		if err != nil {
			t.Fatalf("Không thể tạo file backup app.log: %v", err)
		}
		hasBackup = true
	}

	// Đảm bảo dọn dẹp và khôi phục log file sau khi test xong
	defer func() {
		_ = os.Remove("app.log") // Xóa file log test nếu còn
		if hasBackup {
			_ = os.Rename("app.log.bak", "app.log")
		}
	}()

	s := service.NewHealthService(nil)

	// Case 1: File log không tồn tại
	t.Run("Log file does not exist", func(t *testing.T) {
		_ = os.Remove("app.log") // Chắc chắn là không tồn tại

		logs, err := s.GetAppLogs(5)
		if err != nil {
			t.Fatalf("Không mong đợi lỗi khi file log không tồn tại, nhưng nhận được: %v", err)
		}

		if len(logs) != 1 {
			t.Fatalf("Mong đợi kết quả có 1 phần tử báo lỗi, nhận được %d", len(logs))
		}

		expectedMsg := "Không tìm thấy log file"
		if !strings.Contains(logs[0], expectedMsg) {
			t.Errorf("Mong đợi log chứa %q, nhận được %q", expectedMsg, logs[0])
		}
	})

	// Case 2: File log tồn tại và đọc N dòng cuối cùng
	t.Run("Read last N lines", func(t *testing.T) {
		logContent := "line 1\nline 2\nline 3\nline 4\nline 5\n"
		err := os.WriteFile("app.log", []byte(logContent), 0644)
		if err != nil {
			t.Fatalf("Không thể tạo file log test: %v", err)
		}

		// Đọc 3 dòng cuối cùng
		logs, err := s.GetAppLogs(3)
		if err != nil {
			t.Fatalf("Nhận lỗi không mong muốn khi đọc log: %v", err)
		}

		if len(logs) != 3 {
			t.Fatalf("Mong đợi nhận được 3 dòng log, nhưng nhận được %d", len(logs))
		}

		expectedLines := []string{"line 3", "line 4", "line 5"}
		for i, line := range expectedLines {
			if logs[i] != line {
				t.Errorf("Mong đợi dòng %d là %q, nhưng nhận được %q", i, line, logs[i])
			}
		}

		// Đọc nhiều hơn số dòng hiện có
		logsAll, err := s.GetAppLogs(10)
		if err != nil {
			t.Fatalf("Nhận lỗi không mong muốn khi đọc toàn bộ log: %v", err)
		}

		if len(logsAll) != 5 {
			t.Fatalf("Mong đợi nhận được 5 dòng log, nhưng nhận được %d", len(logsAll))
		}

		expectedAllLines := []string{"line 1", "line 2", "line 3", "line 4", "line 5"}
		for i, line := range expectedAllLines {
			if logsAll[i] != line {
				t.Errorf("Mong đợi dòng %d là %q, nhưng nhận được %q", i, line, logsAll[i])
			}
		}
	})
}
