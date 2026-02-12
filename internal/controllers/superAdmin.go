package controllers

import (
	"easyfind/internal/apiErr"
	"easyfind/internal/models"
	"easyfind/internal/services"
	"easyfind/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Force import models to make swag happy
var _ = models.Account{}

// --- 数据统计 ---

// GetSystemStats godoc
// @Summary 获取系统统计数据
// @Description 超级管理员获取全局统计数据
// @Tags SuperAdmin
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=models.SystemStats} "获取成功"
// @Router /api/v1/super/stats [get]
func GetSystemStats(c *gin.Context) {
	stats, err := services.SuperAdminServiceApp.GetSystemStats()
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}
	response.Success(c, stats)
}

// --- 用户管理 ---

type GetUserListRequest struct {
	Role     int    `form:"role"`
	Keyword  string `form:"keyword"`
	PageNum  int    `form:"page_num"`
	PageSize int    `form:"page_size"`
}

// GetUserList godoc
// @Summary 获取用户列表
// @Description 获取所有用户列表 (支持按角色筛选)
// @Tags SuperAdmin
// @Accept json
// @Produce json
// @Param role query int false "角色 (1:用户, 2:管理员)"
// @Param keyword query string false "关键词 (用户名/姓名/昵称)"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=map[string]interface{}} "获取成功"
// @Router /api/v1/super/users [get]
func GetUserList(c *gin.Context) {
	var req GetUserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	users, total, err := services.SuperAdminServiceApp.GetUserList(req.Role, req.Keyword, req.PageNum, req.PageSize)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	response.Success(c, gin.H{
		"list":  users,
		"total": total,
	})
}

type CreateAdminRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

// AdminCreateUser godoc
// @Summary 创建失物招领管理员
// @Description 创建新的失物招领管理员账号
// @Tags SuperAdmin
// @Accept json
// @Produce json
// @Param request body CreateAdminRequest true "管理员信息"
// @Success 200 {object} response.Response "创建成功"
// @Router /api/v1/super/users/admin [post]
func AdminCreateUser(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	if err := services.SuperAdminServiceApp.CreateAdminUser(req.Username, req.Password, req.Name); err != nil {
		if err.Error() == "用户名已存在" {
			apiErr.HandleBizError(c, response.ErrUserExist, "用户名已存在")
		} else {
			apiErr.HandleSysError(c, response.ErrCreateFail, err)
		}
		return
	}
	response.Success(c, nil)
}

type UpdateUserStatusRequest struct {
	IsActive bool `json:"is_active"`
}

// AdminUpdateUserStatus godoc
// @Summary 禁用/启用账号
// @Description 修改账号激活状态
// @Tags SuperAdmin
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body UpdateUserStatusRequest true "状态"
// @Success 200 {object} response.Response "操作成功"
// @Router /api/v1/super/users/{id}/status [put]
func AdminUpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "ID格式错误")
		return
	}

	var req UpdateUserStatusRequest
	// ShouldBindJSON 对于 bool 字段，如果传 false 它是默认值，所以不能用 binding:"required" 简单判断
	// 这里通过 body 传值
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	if err := services.SuperAdminServiceApp.UpdateUserStatus(uint(id), req.IsActive); err != nil {
		apiErr.HandleSysError(c, response.CodeError, err)
		return
	}
	response.Success(c, nil)
}

// --- 分类管理 ---

type AddCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

// AddCategory godoc
// @Summary 添加物品分类
// @Description 添加新的物品分类
// @Tags SuperAdmin
// @Accept json
// @Produce json
// @Param request body AddCategoryRequest true "分类信息"
// @Success 200 {object} response.Response "添加成功"
// @Router /api/v1/super/categories [post]
func AddCategory(c *gin.Context) {
	var req AddCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	if err := services.SuperAdminServiceApp.AddCategory(req.Name); err != nil {
		apiErr.HandleBizError(c, response.CodeError, err.Error())
		return
	}
	response.Success(c, nil)
}

// DeleteCategory godoc
// @Summary 删除物品分类
// @Description 删除指定的物品分类
// @Tags SuperAdmin
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response "删除成功"
// @Router /api/v1/super/categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "ID格式错误")
		return
	}

	if err := services.SuperAdminServiceApp.DeleteCategory(id); err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}
	response.Success(c, nil)
}

// --- 公告管理 ---

type CreateAnnouncementRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Type    string `json:"type"`   // global or region
	Region  string `json:"region"` // if type is region
	IsTop   bool   `json:"is_top"`
}

// CreateAnnouncement godoc
// @Summary 发布系统公告
// @Description 发布新的系统公告
// @Tags SuperAdmin
// @Accept json
// @Produce json
// @Param request body CreateAnnouncementRequest true "公告内容"
// @Success 200 {object} response.Response "发布成功"
// @Router /api/v1/super/announcements [post]
func CreateAnnouncement(c *gin.Context) {
	var req CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	publisherName := "System Admin"
	if name, exists := c.Get("name"); exists {
		publisherName = name.(string)
	}

	// Super admins can create "published" announcements directly if they want,
	// but the service logic defaults "region" to "pending".
	// Let's assume super admin regional announcements are auto-approved?
	// For now, let's just stick to the service logic. If type is region, it goes pending.
	// If super admin wants it published immediately, they can just approve it or we can add logic.
	// However, usually "region" announcements come from regional admins.
	// If a super admin posts a regional announcement, maybe it should be auto-approved.
	// Let's keep it simple: Service decides status based on type.

	// Actually, let's override logic in service if we want, but the user requirement is:
	// "管理员还可以发布区域公告，然后被系统管理员审核" -> implies Admin (Role 2) posts, Super Admin (Role 3) approves.
	// If Super Admin posts, it should probably just be published.
	// But let's assume Super Admin mainly posts Global.

	if err := services.SuperAdminServiceApp.CreateAnnouncement(req.Title, req.Content, req.Type, req.Region, publisherName, req.IsTop); err != nil {
		apiErr.HandleSysError(c, response.ErrCreateFail, err)
		return
	}
	response.Success(c, nil)
}

// GetAnnouncementsByAdmin godoc
// @Summary 获取公告列表
// @Description 获取系统公告列表 (支持筛选)
// @Tags SuperAdmin
// @Accept json
// @Produce json
// @Param type query string false "类型 (global/region)"
// @Param status query string false "状态 (published/pending/rejected)"
// @Param region query string false "区域"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response "获取成功"
// @Router /api/v1/super/announcements [get]
func GetAnnouncementsByAdmin(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.DefaultQuery("page_num", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	pType := c.Query("type")
	status := c.Query("status")
	region := c.Query("region")

	list, total, err := services.SuperAdminServiceApp.GetAnnouncements(pType, status, region, pageNum, pageSize)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

type ReviewAnnouncementRequest struct {
	ID     uint   `json:"id" binding:"required"`
	Status string `json:"status" binding:"required,oneof=published rejected"` // only allow these two
}

// ReviewAnnouncement godoc
// @Summary 审核公告
// @Description 审核区域公告 (通过/拒绝)
// @Tags SuperAdmin
// @Accept json
// @Produce json
// @Param request body ReviewAnnouncementRequest true "审核信息"
// @Success 200 {object} response.Response "操作成功"
// @Router /api/v1/super/announcements/review [put]
func ReviewAnnouncement(c *gin.Context) {
	var req ReviewAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	if err := services.SuperAdminServiceApp.UpdateAnnouncementStatus(req.ID, req.Status); err != nil {
		if err.Error() == "公告不存在" {
			apiErr.HandleSysError(c, response.ErrTargetNotFound, err)
			return
		}
		apiErr.HandleSysError(c, response.ErrItemUpdateFali, err)
		return
	}
	response.Success(c, nil)
}

// DeleteAnnouncement godoc
// @Summary 删除公告
// @Description 删除指定的公告
// @Tags SuperAdmin
// @Accept json
// @Produce json
// @Param id path int true "公告ID"
// @Success 200 {object} response.Response "删除成功"
// @Router /api/v1/super/announcements/{id} [delete]
func DeleteAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "ID格式错误")
		return
	}

	if err := services.SuperAdminServiceApp.DeleteAnnouncement(uint(id)); err != nil {
		apiErr.HandleSysError(c, response.CodeError, err)
		return
	}
	response.Success(c, nil)
}

// --- 反馈管理 ---

// GetFeedbacks godoc
// @Summary 获取用户反馈
// @Description 获取用户提交的反馈与投诉
// @Tags SuperAdmin
// @Accept json
// @Produce json
// @Param status query string false "状态 (pending, resolved)"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=map[string]interface{}} "获取成功"
// @Router /api/v1/super/feedbacks [get]
func GetFeedbacks(c *gin.Context) {
	status := c.Query("status")
	pageNum, _ := strconv.Atoi(c.DefaultQuery("page_num", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	list, total, err := services.SuperAdminServiceApp.GetFeedbacks(status, pageNum, pageSize)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

type ReplyFeedbackRequest struct {
	Reply string `json:"reply" binding:"required"`
}

// ReplyFeedback godoc
// @Summary 回复反馈
// @Description 管理员回复并处理反馈
// @Tags SuperAdmin
// @Accept json
// @Produce json
// @Param id path int true "反馈ID"
// @Param request body ReplyFeedbackRequest true "回复内容"
// @Success 200 {object} response.Response "操作成功"
// @Router /api/v1/super/feedbacks/{id}/reply [put]
func ReplyFeedback(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "ID格式错误")
		return
	}

	var req ReplyFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	handlerName := "Admin"
	if name, exists := c.Get("name"); exists {
		handlerName = name.(string)
	}

	if err := services.SuperAdminServiceApp.ReplyFeedback(uint(id), req.Reply, handlerName); err != nil {
		apiErr.HandleSysError(c, response.CodeError, err)
		return
	}
	response.Success(c, nil)
}

// --- 数据清理 ---

type CleanupRequest struct {
	Days int `json:"days" binding:"required,min=1"`
}

// CleanupData godoc
// @Summary 清理过期数据
// @Description 清理指定天数前的无效数据(如已取消的帖子)
// @Tags SuperAdmin
// @Accept json
// @Produce json
// @Param request body CleanupRequest true "清理参数"
// @Success 200 {object} response.Response "清理成功"
// @Router /api/v1/super/data/cleanup [post]
func CleanupData(c *gin.Context) {
	var req CleanupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	count, err := services.SuperAdminServiceApp.CleanupData(req.Days)
	if err != nil {
		apiErr.HandleSysError(c, response.CodeError, err)
		return
	}
	response.Success(c, gin.H{"deleted_count": count})
}
