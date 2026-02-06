package controllers

import (
	"easyfind/internal/models"
	"easyfind/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     int    `json:"role" binding:"required"` // 1:学生/老师, 2:失物招领管理员, 3:系统管理员
	Remember bool   `json:"remember"`                // 记住我
}

type LoginResponse struct {
	Token string `json:"token"`
}

// LoginAccount godoc
// @Summary 用户登录
// @Description 用户登录接口，支持学生/老师、失物招领管理员、系统管理员
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录请求参数"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} map[string]string "认证失败"
// @Router /api/v1/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := services.AuthServiceApp.Login(req.Username, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     int    `json:"role" binding:"required"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

// RegisterAccount godoc
// @Summary 注册账号 (测试用)
// @Description 注册新账号，用于测试不同身份
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册请求参数"
// @Success 200 {object} map[string]string "注册成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Router /api/v1/register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	account := models.Account{
		Username: req.Username,
		Password: req.Password, // Service 层会处理加密
		Role:     models.UserRole(req.Role),
		Name:     req.Name,
		Phone:    req.Phone,
	}

	if err := services.AuthServiceApp.Register(account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account created successfully"})
}

// 修改密码
func UpdatePassword(c *gin.Context) {}

// 退出登录
func Logout(c *gin.Context) {}
