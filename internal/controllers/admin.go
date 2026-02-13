package controllers

import (
	"easyfind/internal/apiErr"
	"easyfind/internal/models"
	"easyfind/internal/services"
	"easyfind/pkg/response"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

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

type AdminUpdateRecordRequest struct {
	Title        string `json:"title"`
	Campus       string `json:"campus"`
	Category     string `json:"category"`
	Location     string `json:"location"`
	Time         string `json:"time"`
	Description  string `json:"description"`
	Img1         string `json:"img1"`
	Img2         string `json:"img2"`
	Img3         string `json:"img3"`
	Img4         string `json:"img4"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
}

// AdminUpdateRecord 管理员更新Found类型的物品信息
func AdminUpdateRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "ID格式错误")
		return
	}

	var req AdminUpdateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	oldRecord, err := services.GetRecordById(uint(id))
	if err != nil {
		apiErr.HandleBizError(c, response.ErrItemNotFound, "记录不存在")
		return
	}
	if oldRecord.Type != models.TypeFound {
		apiErr.HandleBizError(c, response.ErrNoPerMission, "只能修改招领信息")
		return
	}

	services.RemoveImg(oldRecord, req.Img1, req.Img2, req.Img3, req.Img4)

	var t time.Time
	if req.Time != "" {
		var err error
		t, err = time.ParseInLocation("2006-01-02 15:04:05", req.Time, time.Local)
		if err != nil {
			apiErr.HandleBizError(c, response.InvalidParams, "时间格式错误")
			return
		}
	}

	item := models.Item{
		Title:        req.Title,
		Campus:       req.Campus,
		Category:     req.Category,
		Location:     req.Location,
		Time:         t,
		Description:  req.Description,
		Img1:         req.Img1,
		Img2:         req.Img2,
		Img3:         req.Img3,
		Img4:         req.Img4,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
	}

	if err := services.UpdateRecord(id, item); err != nil {
		apiErr.HandleSysError(c, response.ErrItemUpdateFali, err)
		return
	}

	response.Success(c, nil)
}

// GetSystemStatsByAdmin 管理员获取系统数据
func GetSystemStatsByAdmin(c *gin.Context) {
	stats, err := services.GetSystemStats()
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}
	response.Success(c, stats)
}

// ExportStats 以纯文本形式导出数据
func ExportStats(c *gin.Context) {
	stats, err := services.GetSystemStats()
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=stats.csv")
	c.Header("Transfer-Encoding", "chunked")

	writer := csv.NewWriter(c.Writer)
	err = writer.Write([]string{"总发布量", "匹配数量", "匹配率", "认领数量", "认领率", "归档数量", "归档率"})
	if err != nil {
		return
	}

	err = writer.Write([]string{
		fmt.Sprintf("%d", stats.Total),
		fmt.Sprintf("%d", stats.Matched),
		fmt.Sprintf("%.2f%%", stats.MatchedRate*100),
		fmt.Sprintf("%d", stats.Claimed),
		fmt.Sprintf("%.2f%%", stats.ClaimedRate*100),
		fmt.Sprintf("%d", stats.Archived),
		fmt.Sprintf("%.2f%%", stats.ArchivedRate*100),
	})
	if err != nil {
		return
	}
	writer.Flush()
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
