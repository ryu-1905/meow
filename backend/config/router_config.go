package config

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter cấu hình Gin Engine, thiết lập CORS và định nghĩa các route API.
func SetupRouter(container *Container) *gin.Engine {
	r := gin.Default()

	// Sử dụng middleware CORS chính thức của Gin để cho phép frontend Svelte kết nối an toàn
	r.Use(cors.Default())

	// API Routes
	api := r.Group("/api")
	{
		api.GET("/health", container.HealthController.GetHealth)

		auth := api.Group("/auth")
		{
			auth.POST("/register", container.UserController.Register)
			auth.POST("/login", container.UserController.Login)
		}
	}

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
