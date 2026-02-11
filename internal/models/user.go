package models

import (
	"gorm.io/gorm"
)

type UserRole int

const (
	RoleStudentTeacher UserRole = 1
	RoleLFAdmin        UserRole = 2
	RoleSysAdmin       UserRole = 3
)

type Account struct {
	gorm.Model
	Username   string   `gorm:"type:varchar(191);uniqueIndex;not null" json:"username"` // 学号或工号
	Password   string   `gorm:"not null" json:"-"`
	Role       UserRole `gorm:"default:1" json:"role"` // 1:学生/老师, 2:失物招领管理员, 3:系统管理员
	Name       string   `json:"name"`
	Nickname   string   `json:"nickname"` // 昵称
	Phone      string   `json:"phone"`
	IsActive   bool     `gorm:"default:true" json:"is_active"`
	FirstLogin bool     `gorm:"default:true" json:"first_login"` // 是否首次登录
}

// TableName overrides the table name used by User to `profiles`
func (Account) TableName() string {
	return "accounts"
}
