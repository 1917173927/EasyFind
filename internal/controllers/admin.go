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

// GetRecordByAdmin 管理员获取帖子列表
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

// GetPendingRecordByAdmin 获取待审核列表
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

// ApproveRecord 通过审核
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

// RejectRecord 驳回审核
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

// ArchiveRecord 归档帖子
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

// GetPendingClaimByAdmin 管理员获取待审核认领列表
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

// ApproveClaim 通过认领申请
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

// RejectClaim 驳回认领申请
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
