package config

import (
	"sync"
	"testing"
)

// ResetDBPoolForTest reset lại sync.Once và dbPool phục vụ cho mục đích kiểm thử.
func ResetDBPoolForTest() {
	once = sync.Once{}
	if dbPool != nil {
		dbPool.Close()
		dbPool = nil
	}
}

func TestGetDB_Initial(t *testing.T) {
	ResetDBPoolForTest()
	db := GetDB()
	if db != nil {
		t.Errorf("Mong đợi GetDB() trả về nil trước khi khởi tạo, nhưng nhận được: %v", db)
	}
}

func TestInitDB_InvalidConnStr(t *testing.T) {
	ResetDBPoolForTest()
	defer ResetDBPoolForTest()

	// Sử dụng connection string không hợp lệ (ví dụ: port không phải số) để kích hoạt lỗi parse config
	connStr := "postgres://user:password@localhost:invalid_port/dbname"
	pool, err := InitDB(connStr)
	if err == nil {
		t.Errorf("Mong đợi InitDB() trả về lỗi với connection string không hợp lệ, nhưng không nhận được lỗi")
	}
	if pool != nil {
		t.Errorf("Mong đợi pool trả về là nil khi lỗi, nhưng nhận được: %v", pool)
	}
}
