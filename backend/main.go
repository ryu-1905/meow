package main

import (
	"github.com/ryu-1905/meow/config"
	_ "github.com/ryu-1905/meow/doc"
)

// @title Meow API
// @version 1.0
// @description This is the backend API for Meow application.
// @schemes https http
// @BasePath /api

func main() {
	// Khởi chạy ứng dụng
	config.ConfigApp()
}
