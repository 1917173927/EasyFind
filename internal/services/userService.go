package services

import (
	"easyfind/internal/models"
	"easyfind/pkg/database"
	"errors"

	"gorm.io/gorm"
)

type UserService struct{}

var UserServiceApp = new(UserService)

// GetUserByUsername 根据用户名获取用户信息
func (s *UserService) GetUserByUsername(username string) (*models.Account, error) {
	var user models.Account
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUserProfile 更新个人信息
func (s *UserService) UpdateUserProfile(id uint, name, nickname, phone, avatar string) error {
	updates := make(map[string]interface{})

	// 只有非空值才更新
	if name != "" {
		updates["name"] = name
	}
	if nickname != "" {
		updates["nickname"] = nickname
	}
	if phone != "" {
		updates["phone"] = phone
	}
	if avatar != "" {
		updates["avatar"] = avatar
	}

	if len(updates) == 0 {
		return nil // 无需更新
	}

	updates["first_login"] = false

	return database.DB.Model(&models.Account{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteUserByUsername 硬删除用户
func (s *UserService) DeleteUserByUsername(username string) error {
	result := database.DB.Unscoped().Where("username = ?", username).Delete(&models.Account{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

// CreateFeedback 创建用户反馈
func (s *UserService) CreateFeedback(userID uint, fType, content, contact string) error {
	feedback := models.Feedback{
		UserID:  userID,
		Type:    fType,
		Content: content,
		Contact: contact,
		Status:  "pending",
	}
	return database.DB.Create(&feedback).Error
}
