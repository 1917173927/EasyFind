package controllers

import (
	"easyfind/internal/apiErr"
	"easyfind/internal/services"
	"easyfind/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetRecordByAdminRequest struct {
	Campus      string `form:"campus"`
	Category    string `form:"category"`
	LostOrFound int    `form:"lost_or_found"`
	Status      string `form:"status"`
	PageNum     int    `form:"page_num"`
	PageSize    int    `form:"page_size"`
}

// GetRecordByAdmin godoc
// @Summary 管理员获取物品列表
// @Description 管理员获取所有物品列表，支持筛选
// @Tags Admin (Item)
// @Accept json
// @Produce json
// @Param campus query string false "校区"
// @Param category query string false "分类"
// @Param lost_or_found query int false "物品类型 (1:Lost, 2:Found)"
// @Param status query string false "状态"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=map[string]interface{}} "获取成功"
// @Router /api/v1/admin/items [get]
func GetRecordByAdmin(c *gin.Context) {
	var req GetRecordByAdminRequest
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

	records, err := services.GetRecordByAdmin(req.Campus, req.Category, req.LostOrFound, req.Status, req.PageNum, req.PageSize)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	total, err := services.GetTotalPageNumByAdmin(req.Campus, req.Category, req.LostOrFound, req.Status)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	response.Success(c, gin.H{
		"list":  records,
		"total": total,
	})
}

type GetPendingRecordByAdminRequest struct {
	LostOrFound int `form:"lost_or_found"`
	PageNum     int `form:"page_num"`
	PageSize    int `form:"page_size"`
}

// GetPendingRecordByAdmin godoc
// @Summary 管理员获取待审核物品
// @Description 管理员获取所有待审核的物品列表
// @Tags Admin (Item)
// @Accept json
// @Produce json
// @Param lost_or_found query int false "物品类型 (1:Lost, 2:Found)"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=map[string]interface{}} "获取成功"
// @Router /api/v1/admin/items/pending [get]
func GetPendingRecordByAdmin(c *gin.Context) {
	var req GetPendingRecordByAdminRequest
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

	records, err := services.GetPendingRecordByAdmin(req.LostOrFound, req.PageNum, req.PageSize)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	total, err := services.GetPendingTotalPageNumByAdmin(req.LostOrFound)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	response.Success(c, gin.H{
		"list":  records,
		"total": total,
	})
}

// ApproveRecord godoc
// @Summary 通过审核
// @Description 管理员通过物品发布的审核
// @Tags Admin (Item)
// @Accept json
// @Produce json
// @Param id path int true "物品ID"
// @Success 200 {object} response.Response "操作成功"
// @Router /api/v1/admin/items/{id}/approve [put]
func ApproveRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "ID格式错误")
		return
	}

	if err := services.ApproveRecord(uint(id)); err != nil {
		apiErr.HandleSysError(c, response.ErrItemUpdateFali, err)
		return
	}

	response.Success(c, nil)
}

type RejectRecordRequest struct {
	RejectReason string `json:"reject_reason" binding:"required"`
}

// RejectRecord godoc
// @Summary 驳回审核
// @Description 管理员驳回物品发布的审核，需提供理由
// @Tags Admin (Item)
// @Accept json
// @Produce json
// @Param id path int true "物品ID"
// @Param request body RejectRecordRequest true "驳回理由"
// @Success 200 {object} response.Response "操作成功"
// @Router /api/v1/admin/items/{id}/reject [put]
func RejectRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "ID格式错误")
		return
	}

	var req RejectRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	if err := services.RejectRecord(uint(id), req.RejectReason); err != nil {
		apiErr.HandleSysError(c, response.ErrItemUpdateFali, err)
		return
	}

	response.Success(c, nil)
}

type ArchiveRecordRequest struct {
	ProcessMethod string `json:"process_method" binding:"required"`
}

// ArchiveRecord godoc
// @Summary 归档帖子
// @Description 管理员归档物品帖子
// @Tags Admin (Item)
// @Accept json
// @Produce json
// @Param id path int true "物品ID"
// @Param request body ArchiveRecordRequest true "处理方式"
// @Success 200 {object} response.Response "归档成功"
// @Router /api/v1/admin/items/{id}/archive [put]
func ArchiveRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "ID格式错误")
		return
	}

	var req ArchiveRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	if err := services.ArchiveRecord(uint(id), req.ProcessMethod); err != nil {
		apiErr.HandleSysError(c, response.ErrItemUpdateFali, err)
		return
	}

	response.Success(c, nil)
}

type GetPendingClaimByAdminRequest struct {
	PageNum  int `form:"page_num"`
	PageSize int `form:"page_size"`
}

// GetPendingClaimByAdmin godoc
// @Summary 管理员获取待审核认领
// @Description 管理员获取所有待审核的认领申请
// @Tags Admin (Claim)
// @Accept json
// @Produce json
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=map[string]interface{}} "获取成功"
// @Router /api/v1/admin/claims/pending [get]
func GetPendingClaimByAdmin(c *gin.Context) {
	var req GetPendingClaimByAdminRequest
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

	claims, err := services.GetPendingClaimByAdmin(req.PageNum, req.PageSize)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	total, err := services.GetPendingClaimTotalPageNumByAdmin()
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	response.Success(c, gin.H{
		"list":  claims,
		"total": total,
	})
}

// ApproveClaim godoc
// @Summary 通过认领申请
// @Description 管理员通过认领申请
// @Tags Admin (Claim)
// @Accept json
// @Produce json
// @Param id path int true "认领ID"
// @Success 200 {object} response.Response "操作成功"
// @Router /api/v1/admin/claims/{id}/approve [put]
func ApproveClaim(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "ID格式错误")
		return
	}

	if err := services.ApproveClaim(uint(id)); err != nil {
		apiErr.HandleSysError(c, response.ErrClaimUpdateFail, err)
		return
	}

	response.Success(c, nil)
}

// RejectClaim godoc
// @Summary 驳回认领申请
// @Description 管理员驳回认领申请
// @Tags Admin (Claim)
// @Accept json
// @Produce json
// @Param id path int true "认领ID"
// @Success 200 {object} response.Response "操作成功"
// @Router /api/v1/admin/claims/{id}/reject [put]
func RejectClaim(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "ID格式错误")
		return
	}

	if err := services.RejectClaim(uint(id)); err != nil {
		apiErr.HandleSysError(c, response.ErrClaimUpdateFail, err)
		return
	}

	response.Success(c, nil)
}
