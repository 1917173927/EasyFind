package controllers

import (
	"easyfind/internal/models"
	"easyfind/internal/services"
	"easyfind/pkg/response"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     int    `json:"role" binding:"required"` // 1:学生/老师, 2:失物招领管理员, 3:系统管理员
	Remember bool   `json:"remember"`                // 记住我
}

// LoginResponse 响应数据结构
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
// @Success 200 {object} response.Response{data=LoginResponse} "登录成功"
// @Failure 200 {object} response.Response "认证失败"
// @Router /api/v1/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.InvalidParams)
		return
	}

	token, err := services.AuthServiceApp.Login(req.Username, req.Password, req.Role)
	if err != nil {
		// 这里简单处理，实际可根据 error 类型返回不同 code
		if err.Error() == "密码错误" {
			response.Fail(c, response.ErrPasswordWrong)
		} else if err.Error() == "用户不存在或角色不匹配" {
			response.Fail(c, response.ErrUserNotExist)
		} else {
			response.FailWithMessage(c, response.CodeError, err.Error())
		}
		return
	}

	// 统一返回格式：code, msg, data(token)
	response.Success(c, LoginResponse{Token: token})
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
// @Success 200 {object} response.Response "注册成功"
// @Failure 200 {object} response.Response "注册失败"
// @Router /api/v1/register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.InvalidParams)
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
		if err.Error() == "用户名已存在" {
			response.Fail(c, response.ErrUserExist)
		} else {
			response.FailWithMessage(c, response.CodeError, err.Error())
		}
		return
	}

	response.Success(c, nil)
}

// 修改密码
func UpdatePassword(c *gin.Context) {}

// 退出登录
func Logout(c *gin.Context) {}
