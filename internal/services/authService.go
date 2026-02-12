package services

import (
	"easyfind/internal/models"
	"easyfind/pkg/database"
	"easyfind/pkg/utils"
	"errors"

	"gorm.io/gorm"
)

type AuthService struct{}

var AuthServiceApp = new(AuthService)

// Login 登录业务逻辑
func (s *AuthService) Login(username, password string, role int) (string, error) {
	var account models.Account
	// 查询用户，根据用户名和角色筛选
	if err := database.DB.Where("username = ? AND role = ?", username, role).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("用户不存在或角色不匹配")
		}
		return "", err
	}

	// 检查账号状态
	if !account.IsActive {
		return "", errors.New("账号已被冻结")
	}

	// 验证密码
	if !utils.CheckPasswordHash(password, account.Password) {
		return "", errors.New("密码错误")
	}

	// 生成 Token
	token, err := utils.GenerateToken(account.ID, account.Username, int(account.Role))
	if err != nil {
		return "", errors.New("生成令牌失败")
	}

	return token, nil
}

// Register 注册业务逻辑
func (s *AuthService) Register(req models.Account) error {
	// 检查用户名是否已存在
	var count int64
	database.DB.Model(&models.Account{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("密码加密失败")
	}
	req.Password = hashedPassword

	// 设置默认头像
	if req.Avatar == "" {
		req.Avatar = "default.jpeg"
	}

	// 创建用户
	if err := database.DB.Create(&req).Error; err != nil {
		return err
	}

	return nil
}

// UpdatePassword 修改密码业务逻辑
func (s *AuthService) UpdatePassword(userID uint, oldPassword, newPassword string) error {
	var account models.Account
	// 查询用户
	if err := database.DB.First(&account, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !utils.CheckPasswordHash(oldPassword, account.Password) {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	if err := database.DB.Model(&account).Update("password", hashedPassword).Error; err != nil {
		return err
	}

	return nil
}
