package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbPool *pgxpool.Pool
	once   sync.Once
)

// InitDB khởi tạo Connection Pool tới PostgreSQL
func InitDB(connStr string) (*pgxpool.Pool, error) {
	var err error
	once.Do(func() {
		// Tạo config từ connection string
		config, parseErr := pgxpool.ParseConfig(connStr)
		if parseErr != nil {
			err = fmt.Errorf("không thể parse database url: %w", parseErr)
			return
		}

		// Tạo pool kết nối
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		dbPool, parseErr = pgxpool.NewWithConfig(ctx, config)
		if parseErr != nil {
			err = fmt.Errorf("không thể kết nối và tạo pool database: %w", parseErr)
			return
		}

		// Ping thử để đảm bảo kết nối hoạt động
		if pingErr := dbPool.Ping(ctx); pingErr != nil {
			err = fmt.Errorf("không thể ping database: %w", pingErr)
			// Đóng pool nếu lỗi ping xảy ra trong quá trình khởi tạo để dọn dẹp tài nguyên
			dbPool.Close()
			dbPool = nil
			return
		}

		log.Println("Kết nối cơ sở dữ liệu PostgreSQL đã được thiết lập thành công qua Connection Pool.")
	})

	return dbPool, err
}

// GetDB trả về Pool kết nối hiện tại
func GetDB() *pgxpool.Pool {
	return dbPool
}
