package config

import (
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ryu-1905/meow/controller"
	"github.com/ryu-1905/meow/repository"
	"github.com/ryu-1905/meow/service"
)

// Container lưu trữ tất cả các dependencies của ứng dụng sau khi được tiêm (inject).
type Container struct {
	DBPool           *pgxpool.Pool
	HealthService    *service.HealthService
	HealthController *controller.HealthController
	JWTService       *service.JWTService
	UserRepository   repository.UserRepository
	UserService      *service.UserService
	UserController   *controller.UserController
}

// NewContainer khởi tạo và tiêm (inject) các dependencies cần thiết cho ứng dụng.
func NewContainer(dbPool *pgxpool.Pool) *Container {
	healthService := service.NewHealthService(dbPool)
	healthController := controller.NewHealthController(healthService)
	userRepository := repository.NewUserRepository(dbPool)

	// Đọc cấu hình JWT từ biến môi trường
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "019ec244-d465-7ed8-9623-b58a28ac95aa"
		slog.Warn("JWT_SECRET không được thiết lập, sử dụng khoá mặc định dùng cho development")
	}

	jwtService := service.NewJWTService(jwtSecret)
	userService := service.NewUserService(userRepository, jwtService)
	userController := controller.NewUserController(userService)

	return &Container{
		DBPool:           dbPool,
		HealthService:    healthService,
		HealthController: healthController,
		JWTService:       jwtService,
		UserRepository:   userRepository,
		UserService:      userService,
		UserController:   userController,
	}
}
