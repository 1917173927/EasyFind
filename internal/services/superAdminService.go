package services

import (
	"easyfind/internal/models"
	"easyfind/internal/ws"
	"easyfind/pkg/database"
	"easyfind/pkg/utils"
	"errors"
	"time"
)

type SuperAdminService struct{}

type UserListItem struct {
	ID          uint       `json:"ID"`
	Username    string     `json:"username"`
	Role        int        `json:"role"`
	RoleName    string     `json:"role_name"`
	Name        string     `json:"name"`
	Nickname    string     `json:"nickname"`
	Phone       string     `json:"phone"`
	Avatar      string     `json:"avatar"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"CreatedAt"`
	LastLoginAt *time.Time `json:"last_login_at"`
}

var SuperAdminServiceApp = new(SuperAdminService)

// --- 统计数据 ---

func (s *SuperAdminService) GetSystemStats() (*models.SystemStats, error) {
	var stats models.SystemStats

	database.DB.Model(&models.Account{}).Count(&stats.TotalUsers)

	// 活跃用户：最近30天登录过？目前模型没有 LastLogin 字段，暂时统计 isActive 的用户
	database.DB.Model(&models.Account{}).Where("is_active = ?", true).Count(&stats.ActiveUsers)

	database.DB.Model(&models.Item{}).Count(&stats.TotalItems)

	database.DB.Model(&models.Item{}).
		Where("status IN ?", []models.ItemStatus{models.StatusMatched, models.StatusClaimed}).
		Count(&stats.SolvedItems)

	database.DB.Model(&models.Claim{}).Count(&stats.TotalClaims)

	// 今日新增
	todayStart := time.Now().Truncate(24 * time.Hour)
	database.DB.Model(&models.Item{}).Where("created_at >= ?", todayStart).Count(&stats.TodayItems)

	return &stats, nil
}

// --- 用户管理 ---

func (s *SuperAdminService) GetUserList(role int, keyword string, pageNum, pageSize int) ([]UserListItem, int64, error) {
	var users []models.Account
	var total int64

	query := database.DB.Model(&models.Account{})

	if role > 0 {
		query = query.Where("role = ?", role)
	}
	if keyword != "" {
		query = query.Where("username LIKE ? OR name LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Select("id, username, role, name, nickname, phone, is_active, created_at, avatar, last_login_at").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Order("created_at desc").
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	list := make([]UserListItem, 0, len(users))
	for _, user := range users {
		roleName := "学生/老师"
		switch user.Role {
		case models.RoleLFAdmin:
			roleName = "管理员"
		case models.RoleSysAdmin:
			roleName = "超级管理员"
		}

		list = append(list, UserListItem{
			ID:          user.ID,
			Username:    user.Username,
			Role:        int(user.Role),
			RoleName:    roleName,
			Name:        user.Name,
			Nickname:    user.Nickname,
			Phone:       user.Phone,
			Avatar:      user.Avatar,
			IsActive:    user.IsActive,
			CreatedAt:   user.CreatedAt,
			LastLoginAt: user.LastLoginAt,
		})
	}

	return list, total, nil
}

func (s *SuperAdminService) UpdateUserStatus(id uint, isActive bool) error {
	return database.DB.Model(&models.Account{}).Where("id = ?", id).Update("is_active", isActive).Error
}

func (s *SuperAdminService) CreateAdminUser(username, password, name string) error {
	// 检查是否存在
	var count int64
	database.DB.Model(&models.Account{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	admin := models.Account{
		Username: username,
		Password: hashedPassword,
		Role:     models.RoleLFAdmin, // 2: 失物招领管理员
		Name:     name,
		IsActive: true,
	}

	return database.DB.Create(&admin).Error
}

// --- 分类管理 ---

func (s *SuperAdminService) AddCategory(name string) error {
	var cat models.LostCategory
	if err := database.DB.Where("category_name = ?", name).First(&cat).Error; err == nil {
		return errors.New("分类已存在")
	}
	newCat := models.LostCategory{CategoryName: name}
	return database.DB.Create(&newCat).Error
}

func (s *SuperAdminService) DeleteCategory(id int) error {
	return database.DB.Delete(&models.LostCategory{}, id).Error
}

// --- 公告管理 ---

func (s *SuperAdminService) CreateAnnouncement(title, content, pType, region, publisher string, isTop bool) error {
	status := "published"
	if pType == "region" {
		status = "pending"
	}

	announce := models.Announcement{
		Title:     title,
		Content:   content,
		Type:      pType,
		Region:    region,
		Publisher: publisher,
		IsTop:     isTop,
		Status:    status,
	}
	if err := database.DB.Create(&announce).Error; err != nil {
		return err
	}
	if status == "published" {
		ws.BroadcastUpdate(ws.ScopeAnnouncement, announce.ID)
	}
	return nil
}

func (s *SuperAdminService) GetAnnouncements(pType, status, region string, page, size int) ([]models.Announcement, int64, error) {
	var list []models.Announcement
	var total int64

	db := database.DB.Model(&models.Announcement{})
	if pType != "" {
		db = db.Where("type = ?", pType)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if region != "" {
		db = db.Where("region LIKE ?", "%"+region+"%")
	}

	db.Count(&total)

	err := db.Order("is_top desc, created_at desc").
		Limit(size).Offset((page - 1) * size).
		Find(&list).Error

	return list, total, err
}

func (s *SuperAdminService) UpdateAnnouncementStatus(id uint, status string) error {
	res := database.DB.Model(&models.Announcement{}).Where("id = ?", id).Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("公告不存在")
	}
	if status == "published" {
		ws.BroadcastUpdate(ws.ScopeAnnouncement, id)
	}
	return nil
}

func (s *SuperAdminService) DeleteAnnouncement(id uint) error {
	return database.DB.Delete(&models.Announcement{}, id).Error
}

// --- 反馈管理 ---

func (s *SuperAdminService) GetFeedbacks(status string, page, size int) ([]models.Feedback, int64, error) {
	var list []models.Feedback
	var total int64

	query := database.DB.Model(&models.Feedback{}).Preload("User")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Order("created_at desc").Limit(size).Offset((page - 1) * size).Find(&list).Error

	return list, total, err
}

func (s *SuperAdminService) ReplyFeedback(id uint, reply, handler string) error {
	return database.DB.Model(&models.Feedback{}).Where("id = ?", id).Updates(map[string]interface{}{
		"reply":   reply,
		"handler": handler,
		"status":  "resolved",
	}).Error
}

// --- 数据清理 ---

func (s *SuperAdminService) CleanupData(days int) (int64, error) {
	// 清理 days 天前的软删除记录 (真正物理删除)
	// 或者清理 "Cancelled" 且超过 days 天的记录

	deadline := time.Now().AddDate(0, 0, -days)

	// 这里演示清理已取消且超时的 Item
	result := database.DB.Unscoped().Where("status = ? AND updated_at < ?", models.StatusCancelled, deadline).Delete(&models.Item{})

	return result.RowsAffected, result.Error
}
