package controllers

import (
	"easyfind/internal/apiErr"
	"easyfind/internal/services"
	"easyfind/pkg/response"

	"github.com/gin-gonic/gin"
)

type GetClaimProgressRequest struct {
	PageNum  int `form:"page_num"`
	PageSize int `form:"page_size"`
}

// GetClaimProgress godoc
// @Summary 获取招领进度
// @Description 获取当前用户的招领进度（按最新时间倒序）。type: claim_pending|claim_approved|claim_rejected|claim_completed；claim_status: pending|approved|rejected|archived
// @Tags Claim (User)
// @Accept json
// @Produce json
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=map[string]interface{}} "获取成功"
// @Router /api/v1/claims/progress [get]
func GetClaimProgress(c *gin.Context) {
	var req GetClaimProgressRequest
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

	userID, exists := c.Get("userID")
	if !exists {
		apiErr.HandleBizError(c, response.ErrTokenInvalid, "未获取到用户信息")
		return
	}

	list, total, err := services.GetClaimProgress(userID.(uint), req.PageNum, req.PageSize)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}
