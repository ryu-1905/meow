package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ryu-1905/meow/service"
)

// RegisterRequest đại diện cho body yêu cầu đăng ký tài khoản mới.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=4"`
}

// LoginRequest đại diện cho body yêu cầu đăng nhập.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse đại diện cho dữ liệu trả về sau khi xác thực thành công.
type AuthResponse struct {
	Token string `json:"token"`
}

// UserController quản lý các HTTP request liên quan đến người dùng (đăng ký, đăng nhập).
type UserController struct {
	userService *service.UserService
}

// NewUserController khởi tạo một instance mới của UserController.
func NewUserController(us *service.UserService) *UserController {
	return &UserController{
		userService: us,
	}
}

// Register godoc
// @Summary Đăng ký tài khoản người dùng mới
// @Description Tạo tài khoản người dùng mới từ name, email, password và sinh Access Token để đăng nhập ngay lập tức.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Thông tin đăng ký"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]string "Lỗi validation hoặc dữ liệu đầu vào không hợp lệ"
// @Failure 500 {object} map[string]string "Lỗi hệ thống hoặc lỗi lưu cơ sở dữ liệu"
// @Router /auth/register [post]
func (uc *UserController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu đầu vào không hợp lệ: " + err.Error()})
		return
	}

	token, err := uc.userService.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{Token: token})
}

// Login godoc
// @Summary Đăng nhập tài khoản
// @Description Xác thực tài khoản bằng email, password và trả về JWT Access Token.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Thông tin đăng nhập"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]string "Lỗi dữ liệu đầu vào không hợp lệ"
// @Failure 401 {object} map[string]string "Sai tài khoản hoặc mật khẩu"
// @Failure 500 {object} map[string]string "Lỗi hệ thống"
// @Router /auth/login [post]
func (uc *UserController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu đầu vào không hợp lệ: " + err.Error()})
		return
	}

	token, err := uc.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// Trả về Unauthorized nếu sai thông tin đăng nhập hoặc không tìm thấy user
		if err.Error() == "invalid email or password" || strings.Contains(err.Error(), "user not found") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Sai địa chỉ email hoặc mật khẩu"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{Token: token})
}
