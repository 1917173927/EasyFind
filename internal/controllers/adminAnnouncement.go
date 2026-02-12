package controllers

import (
	"easyfind/internal/apiErr"
	"easyfind/internal/services"
	"easyfind/pkg/response"

	"github.com/gin-gonic/gin"
)

// --- 区域公告管理 (Role 2) ---

type CreateRegionalAnnouncementRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Region  string `json:"region" binding:"required"` // Region is required for regional announcements
	IsTop   bool   `json:"is_top"`
}

// CreateRegionalAnnouncement godoc
// @Summary 发布区域公告 (需审核)
// @Description 管理员发布区域公告，状态默认为 pending
// @Tags Admin (Announcement)
// @Accept json
// @Produce json
// @Param request body CreateRegionalAnnouncementRequest true "公告内容"
// @Success 200 {object} response.Response "提交成功，等待审核"
// @Router /api/v1/admin/announcements [post]
func CreateRegionalAnnouncement(c *gin.Context) {
	var req CreateRegionalAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	publisherName := "Admin"
	if name, exists := c.Get("name"); exists {
		publisherName = name.(string)
	}

	// Force Type="region"
	if err := services.SuperAdminServiceApp.CreateAnnouncement(req.Title, req.Content, "region", req.Region, publisherName, req.IsTop); err != nil {
		apiErr.HandleSysError(c, response.ErrCreateFail, err)
		return
	}
	response.Success(c, nil)
}
