package controllers

import (
	"easyfind/internal/apiErr"
	"easyfind/internal/models"
	"easyfind/internal/services"
	"easyfind/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetAllItemRequest struct {
	Campus   string `form:"campus"`
	Category string `form:"category"`
	PageNum  int    `form:"page_num"`
	PageSize int    `form:"page_size"`
}

// GetAllItem 获取所有物品列表 (不区分 Lost/Found)
func GetAllItem(c *gin.Context) {
	var req GetAllItemRequest
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

	records, err := services.GetAllItem(req.Campus, req.Category, req.PageNum, req.PageSize)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	total, err := services.GetAllLostAndFoundTotalPageNum(req.Campus, req.Category)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	response.Success(c, gin.H{
		"list":  records,
		"total": total,
	})
}

type GetRecordRequest struct {
	Campus      string `form:"campus"`
	Category    string `form:"category"`
	LostOrFound int    `form:"lost_or_found"`
	PageNum     int    `form:"page_num"`
	PageSize    int    `form:"page_size"`
}

// GetRecord 获取失物招领列表 (区分 Lost/Found)
func GetRecord(c *gin.Context) {
	var req GetRecordRequest
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

	records, err := services.GetRecord(req.Campus, req.Category, req.LostOrFound, req.PageNum, req.PageSize)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	total, err := services.GetTotalPageNum(req.Campus, req.Category, req.LostOrFound)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	response.Success(c, gin.H{
		"list":  records,
		"total": total,
	})
}

// GetRecordById 获取详情
func GetRecordById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "ID格式错误")
		return
	}

	record, err := services.GetRecordById(uint(id))
	if err != nil {
		apiErr.HandleBizError(c, response.ErrItemNotFound, "记录不存在")
		return
	}

	response.Success(c, record)
}

type CreateRecordRequest struct {
	Title        string `json:"title" binding:"required"`
	Type         string `json:"type" binding:"required,one of=lost found"`
	Campus       string `json:"campus"`
	Category     string `json:"category"`
	Location     string `json:"location"`
	Time         string `json:"time"` // 前端传字符串，后端解析
	Description  string `json:"description"`
	Img1         string `json:"img1"`
	Img2         string `json:"img2"`
	Img3         string `json:"img3"`
	Img4         string `json:"img4"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	IsBounty     bool   `json:"is_bounty"`
}

// CreateRecord 发布失物招领
func CreateRecord(c *gin.Context) {
	var req CreateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		apiErr.HandleBizError(c, response.ErrTokenInvalid, "未获取到用户信息")
		return
	}

	item := models.Item{
		Title:        req.Title,
		Type:         models.ItemType(req.Type),
		Campus:       req.Campus,
		Category:     req.Category,
		Location:     req.Location,
		Description:  req.Description,
		Img1:         req.Img1,
		Img2:         req.Img2,
		Img3:         req.Img3,
		Img4:         req.Img4,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		IsBounty:     req.IsBounty,
		PublisherID:  userID.(uint),
		Status:       models.StatusPending,
	}

	if err := services.CreateRecord(item); err != nil {
		apiErr.HandleSysError(c, response.ErrItemCreateFail, err)
		return
	}

	response.Success(c, nil)
}

type UpdateRecordRequest struct {
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

// UpdateRecord 更新记录
func UpdateRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "ID格式错误")
		return
	}

	var req UpdateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	oldRecord, err := services.GetRecordById(uint(id))
	if err != nil {
		apiErr.HandleBizError(c, response.ErrItemNotFound, "记录不存在")
		return
	}

	userID, exists := c.Get("userID")
	if !exists || oldRecord.PublisherID != userID.(uint) {
		apiErr.HandleBizError(c, response.ErrNoPerMission, "无权修改此记录")
		return
	}

	services.RemoveImg(oldRecord, req.Img1, req.Img2, req.Img3, req.Img4)

	item := models.Item{
		Title:        req.Title,
		Campus:       req.Campus,
		Category:     req.Category,
		Location:     req.Location,
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

// DeleteRecord 删除记录
func DeleteRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "ID格式错误")
		return
	}

	record, err := services.GetRecordById(uint(id))
	if err != nil {
		apiErr.HandleBizError(c, response.ErrItemNotFound, "记录不存在")
		return
	}

	userID, exists := c.Get("userID")
	if !exists || record.PublisherID != userID.(uint) {
		apiErr.HandleBizError(c, response.ErrNoPerMission, "无权删除此记录")
		return
	}

	if err := services.DeleteRecord(uint(id)); err != nil {
		apiErr.HandleSysError(c, response.ErrItemDeleteFail, err)
		return
	}

	response.Success(c, nil)
}

// CancelRecord 取消发布
func CancelRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "ID格式错误")
		return
	}

	record, err := services.GetRecordById(uint(id))
	if err != nil {
		apiErr.HandleBizError(c, response.ErrItemNotFound, "记录不存在")
		return
	}

	userID, exists := c.Get("userID")
	if !exists || record.PublisherID != userID.(uint) {
		apiErr.HandleBizError(c, response.ErrNoPerMission, "无权取消此记录")
		return
	}

	if err := services.CancelRecord(uint(id)); err != nil {
		apiErr.HandleSysError(c, response.ErrItemUpdateFali, err)
		return
	}

	response.Success(c, nil)
}

type GetAllMyRecordRequest struct {
	Status   string `form:"status"`
	PageNum  int    `form:"page_num"`
	PageSize int    `form:"page_size"`
}

// GetAllMyRecord 获取我的发布
func GetAllMyRecord(c *gin.Context) {
	var req GetAllMyRecordRequest
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

	records, err := services.GetAllMyRecord(userID.(uint), models.ItemStatus(req.Status), req.PageNum, req.PageSize)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	total, err := services.GetMyTotalPageNum(userID.(uint), models.ItemStatus(req.Status))
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	response.Success(c, gin.H{
		"list":  records,
		"total": total,
	})
}

// GetCategoryList 获取分类列表
func GetCategoryList(c *gin.Context) {
	categories, err := services.GetCategoryList()
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	response.Success(c, categories)
}
