package controllers

import (
	"easyfind/internal/apiErr"
	"easyfind/internal/models"
	"easyfind/internal/services"
	"easyfind/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type GetAllItemRequest struct {
	Campus    string `form:"campus"`
	Category  string `form:"category"`
	Location  string `form:"location"`
	Days      int    `form:"days"`
	Status    string `form:"status"`
	HasBounty *bool  `form:"has_bounty"`
	PageNum   int    `form:"page_num"`
	PageSize  int    `form:"page_size"`
}

// GetAllItem godoc
// @Summary 获取所有物品列表
// @Description 获取所有物品列表 (不区分 Lost/Found)
// @Tags Item (Public)
// @Accept json
// @Produce json
// @Param campus query string false "校区"
// @Param category query string false "分类"
// @Param location query string false "地点 (模糊搜索)"
// @Param days query int false "最近几天 (例如 3 表示最近3天)"
// @Param status query string false "状态 (默认 approved)"
// @Param has_bounty query bool false "是否有悬赏 (true/false)"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=map[string]interface{}} "获取成功"
// @Router /api/v1/items/all [get]
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

	records, err := services.GetAllItem(req.Campus, req.Category, req.Location, req.Days, req.Status, req.HasBounty, req.PageNum, req.PageSize)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	total, err := services.GetAllLostAndFoundTotalPageNum(req.Campus, req.Category, req.Location, req.Days, req.Status, req.HasBounty)
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
	Location    string `form:"location"`
	Days        int    `form:"days"`
	Status      string `form:"status"`
	HasBounty   *bool  `form:"has_bounty"`
	PageNum     int    `form:"page_num"`
	PageSize    int    `form:"page_size"`
}

// GetRecord godoc
// @Summary 获取失物招领列表
// @Description 获取失物招领列表 (区分 Lost/Found)
// @Tags Item (Public)
// @Accept json
// @Produce json
// @Param campus query string false "校区"
// @Param category query string false "分类"
// @Param lost_or_found query int false "物品类型 (1:Lost, 2:Found)"
// @Param location query string false "地点 (模糊搜索)"
// @Param days query int false "最近几天 (例如 3 表示最近3天)"
// @Param status query string false "状态 (默认 approved)"
// @Param has_bounty query bool false "是否有悬赏 (true/false)"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=map[string]interface{}} "获取成功"
// @Router /api/v1/items [get]
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

	records, err := services.GetRecord(req.Campus, req.Category, req.LostOrFound, req.Location, req.Days, req.Status, req.HasBounty, req.PageNum, req.PageSize)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	total, err := services.GetTotalPageNum(req.Campus, req.Category, req.LostOrFound, req.Location, req.Days, req.Status, req.HasBounty)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	response.Success(c, gin.H{
		"list":  records,
		"total": total,
	})
}

// GetRecordById godoc
// @Summary 获取物品详情
// @Description 根据ID获取物品详情
// @Tags Item (Public)
// @Accept json
// @Produce json
// @Param id path int true "物品ID"
// @Success 200 {object} response.Response{data=models.Item} "获取成功"
// @Router /api/v1/items/{id} [get]
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
	Title        string  `json:"title" binding:"required"`
	Type         string  `json:"type" binding:"required,oneof=lost found"`
	Campus       string  `json:"campus"`
	Category     string  `json:"category"`
	Location     string  `json:"location"`
	Time         string  `json:"time"` // 前端传字符串，后端解析
	Description  string  `json:"description"`
	Img1         string  `json:"img1"`
	Img2         string  `json:"img2"`
	Img3         string  `json:"img3"`
	Img4         string  `json:"img4"`
	ContactName  string  `json:"contact_name"`
	ContactPhone string  `json:"contact_phone"`
	Bounty       float64 `json:"bounty"`
}

// CreateRecord godoc
// @Summary 发布失物招领
// @Description 用户发布失物或招领信息
// @Tags Item (User)
// @Accept json
// @Produce json
// @Param request body CreateRecordRequest true "发布信息"
// @Success 200 {object} response.Response "发布成功"
// @Router /api/v1/items [post]
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

	var t time.Time
	if req.Time != "" {
		var err error
		t, err = time.ParseInLocation("2006-01-02 15:04:05", req.Time, time.Local)
		if err != nil {
			apiErr.HandleBizError(c, response.InvalidParams, "时间格式错误, 需为 2006-01-02 15:04:05")
			return
		}
	} else {
		t = time.Now()
	}

	item := models.Item{
		Title:        req.Title,
		Type:         models.ItemType(req.Type),
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
		Bounty:       req.Bounty,
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

// UpdateRecord godoc
// @Summary 更新记录
// @Description 用户更新自己发布的记录
// @Tags Item (User)
// @Accept json
// @Produce json
// @Param id path int true "物品ID"
// @Param request body UpdateRecordRequest true "更新信息"
// @Success 200 {object} response.Response "更新成功"
// @Router /api/v1/items/{id} [put]
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

	var t time.Time
	if req.Time != "" {
		var err error
		t, err = time.ParseInLocation("2006-01-02 15:04:05", req.Time, time.Local)
		if err != nil {
			apiErr.HandleBizError(c, response.InvalidParams, "时间格式错误, 需为 2006-01-02 15:04:05")
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

// DeleteRecord godoc
// @Summary 删除记录
// @Description 用户删除自己发布的记录
// @Tags Item (User)
// @Accept json
// @Produce json
// @Param id path int true "物品ID"
// @Success 200 {object} response.Response "删除成功"
// @Router /api/v1/items/{id} [delete]
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

// CancelRecord godoc
// @Summary 取消发布
// @Description 用户取消自己发布的记录 (通常用于已找到或不再需要)
// @Tags Item (User)
// @Accept json
// @Produce json
// @Param id path int true "物品ID"
// @Success 200 {object} response.Response "取消成功"
// @Router /api/v1/items/{id}/cancel [put]
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

// GetAllMyRecord godoc
// @Summary 获取我的发布
// @Description 获取当前用户发布的所有记录
// @Tags Item (User)
// @Accept json
// @Produce json
// @Param status query string false "状态"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=map[string]interface{}} "获取成功"
// @Router /api/v1/my/items [get]
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

// GetCategoryList godoc
// @Summary 获取分类列表
// @Description 获取所有物品分类
// @Tags Item (Public)
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]string} "获取成功"
// @Router /api/v1/kinds [get]
func GetCategoryList(c *gin.Context) {
	categories, err := services.GetCategoryList()
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	response.Success(c, categories)
}
