package controllers

import (
	"easyfind/internal/services"
	"easyfind/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetProfile godoc
// @Summary 获取个人资料
// @Description 根据 Username 获取用户资料
// @Tags User
// @Accept json
// @Produce json
// @Param username query string true "用户名"
// @Success 200 {object} response.Response{data=models.Account} "获取成功"
// @Router /api/v1/user/profile [get]
func GetProfile(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		response.FailWithMessage(c, response.INVALID_PARAMS, "username is required")
		return
	}

	user, err := services.UserServiceApp.GetUserByUsername(username)
	if err != nil {
		response.FailWithMessage(c, response.ERROR_USER_NOT_EXIST, err.Error())
		return
	}

	response.Success(c, user)
}

type UpdateProfileRequest struct {
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

// UpdateProfile godoc
// @Summary 修改个人资料
// @Description 修改当前登录用户的资料，为空的字段不修改
// @Tags User
// @Accept json
// @Produce json
// @Param request body UpdateProfileRequest true "需要修改的参数"
// @Success 200 {object} response.Response "修改成功"
// @Router /api/v1/user/profile [put]
func UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.INVALID_PARAMS)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Fail(c, response.ERROR)
		return
	}

	if err := services.UserServiceApp.UpdateUserProfile(userID.(uint), req.Name, req.Nickname, req.Phone); err != nil {
		response.FailWithMessage(c, response.ERROR, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteAccount godoc
// @Summary 注销/删除账号 (硬删除)
// @Description 硬删除账号。普通用户只能注销自己，管理员/系统管理员可以注销任意账号。
// @Tags User
// @Accept json
// @Produce json
// @Param username query string false "目标用户名 (管理员使用)"
// @Success 200 {object} response.Response "注销成功"
// @Router /api/v1/user/account [delete]
func DeleteAccount(c *gin.Context) {
	targetUsername := c.Query("username")

	// 获取当前登录用户信息 (由中间件设置)
	currentUsername, _ := c.Get("username")
	role, _ := c.Get("role")

	// 如果没有指定目标用户，默认为操作自己
	if targetUsername == "" {
		targetUsername = currentUsername.(string)
	}

	// 鉴权逻辑
	isSelf := targetUsername == currentUsername.(string)
	isAdmin := false
	if r, ok := role.(int); ok {
		// 检查是否为管理员 (2:失物招领管理员, 3:系统管理员)
		// 这里假设 2 和 3 都有权删除用户，具体业务可调整
		isAdmin = (r == 2 || r == 3)
	}

	if !isSelf && !isAdmin {
		response.FailWithMessage(c, response.ERROR_AUTH_CHECK_TOKEN_FAIL, "无权操作其他用户的账号")
		return
	}

	// 执行删除
	if err := services.UserServiceApp.DeleteUserByUsername(targetUsername); err != nil {
		response.FailWithMessage(c, response.ERROR, err.Error())
		return
	}

	response.Success(c, nil)
}
