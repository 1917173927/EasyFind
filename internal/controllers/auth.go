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
	Nickname string `json:"nickname"`
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
		Nickname: req.Nickname,
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

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// UpdatePassword godoc
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body UpdatePasswordRequest true "修改密码参数"
// @Success 200 {object} response.Response "修改成功"
// @Failure 200 {object} response.Response "修改失败"
// @Router /api/v1/user/password [put]
func UpdatePassword(c *gin.Context) {
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.INVALID_PARAMS)
		return
	}

	// 从 Context 获取 userID (由 JWT 中间件设置)
	userID, exists := c.Get("userID")
	if !exists {
		// 理论上经过中间件不应该走到这里
		response.Fail(c, response.ERROR)
		return
	}

	if err := services.AuthServiceApp.UpdatePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		response.FailWithMessage(c, response.ERROR, err.Error())
		return
	}

	response.Success(c, nil)
}

// Logout godoc
// @Summary 退出登录
// @Description 退出登录，客户端需自行丢弃 Token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "退出成功"
// @Router /api/v1/logout [post]
func Logout(c *gin.Context) {
	// JWT 是无状态的，服务端默认不存储 Token 状态。
	// 真正的注销需要客户端丢弃 Token，或者服务端维护一个 Token 黑名单 (Redis)。
	// 这里目前仅返回成功响应，提示前端清除 Token。
	response.Success(c, nil)
}
