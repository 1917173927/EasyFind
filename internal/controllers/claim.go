package controllers

import (
	"easyfind/internal/apiErr"
	"easyfind/internal/models"
	"easyfind/internal/services"
	"easyfind/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreatClaimRequest struct {
	ItemID uint   `json:"item_id" binding:"required"`
	Proof  string `json:"proof"`
}

// CreatClaim 申请认领
func CreatClaim(c *gin.Context) {
	var req CreatClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		apiErr.HandleBizError(c, response.ErrTokenInvalid, "未获取到用户信息")
		return
	}

	claim := models.Claim{
		ItemID:     req.ItemID,
		ClaimantID: userID.(uint),
		Proof:      req.Proof,
		Status:     string(models.StatusPending),
	}

	if err := services.CreatClaim(claim); err != nil {
		apiErr.HandleSysError(c, response.ErrClaimCreateFail, err)
		return
	}

	response.Success(c, nil)
}

type GetMyClaimRequest struct {
	PageNum  int `form:"page_num"`
	PageSize int `form:"page_size"`
}

// GetMyClaim 获取我的认领列表
func GetMyClaim(c *gin.Context) {
	var req GetMyClaimRequest
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

	claims, err := services.GetMyClaim(userID.(uint), req.PageNum, req.PageSize)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	total, err := services.GetMyClaimTotalPageNum(userID.(uint))
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	response.Success(c, gin.H{
		"list":  claims,
		"total": total,
	})
}

// GetClaimByID 获取认领详情
func GetClaimByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "ID格式错误")
		return
	}

	claim, err := services.GetClaimByID(uint(id))
	if err != nil {
		apiErr.HandleBizError(c, response.ErrClaimNotFound, "认领记录不存在")
		return
	}

	userID, _ := c.Get("userID")
	if claim.ClaimantID != userID.(uint) {
		item, _ := services.GetRecordById(claim.ItemID)
		if item.PublisherID != userID.(uint) {
			apiErr.HandleBizError(c, response.ErrNoPerMission, "无权查看此记录")
			return
		}
	}

	response.Success(c, claim)
}

// ConfirmClaim 确认认领 (归档)
func ConfirmClaim(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "ID格式错误")
		return
	}

	claim, err := services.GetClaimByID(uint(id))
	if err != nil {
		apiErr.HandleBizError(c, response.ErrClaimNotFound, "认领记录不存在")
		return
	}

	item, err := services.GetRecordById(claim.ItemID)
	if err != nil {
		apiErr.HandleBizError(c, response.ErrItemNotFound, "关联物品不存在")
		return
	}

	userID, _ := c.Get("userID")
	if item.PublisherID != userID.(uint) {
		apiErr.HandleBizError(c, response.ErrNoPerMission, "只有发布者可以确认认领")
		return
	}

	if err := services.ConfirmClaim(uint(id)); err != nil {
		apiErr.HandleSysError(c, response.ErrClaimUpdateFail, err)
		return
	}

	response.Success(c, nil)
}
