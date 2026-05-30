package main

import (
	"io"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ryu-1905/meow/controllers"
	_ "github.com/ryu-1905/meow/docs"
	"github.com/ryu-1905/meow/services"
)

// @title Meow API
// @version 1.0
// @description This is the backend API for Meow application.

// @host localhost:8080
// @BasePath /api

func main() {
	// Cấu hình ghi log kép (terminal + file app.log)
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Không thể khởi tạo file log: %v", err)
	} else {
		// Thiết lập gin writer và log writer để ghi ra cả terminal và file log
		gin.DefaultWriter = io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	r := gin.Default()

	// Sử dụng middleware CORS chính thức của Gin để cho phép frontend Svelte kết nối an toàn
	r.Use(cors.Default())

	// Khởi tạo các services và controllers
	healthService := services.NewHealthService()
	healthController := controllers.NewHealthController(healthService)

	// API Routes
	api := r.Group("/api")
	{
		api.GET("/health", healthController.GetHealth)
	}

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Get PORT from env, fallback to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
