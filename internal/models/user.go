package models

import (
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleStudentTeacher UserRole = "student_teacher"
	RoleLFAdmin        UserRole = "lf_admin"
	RoleSysAdmin       UserRole = "sys_admin"
)

type User struct {
	gorm.Model
	Username   string   `gorm:"uniqueIndex;not null" json:"username"` // 学号或工号
	Password   string   `gorm:"not null" json:"-"`
	Role       UserRole `gorm:"type:enum('student_teacher', 'lf_admin', 'sys_admin');default:'student_teacher'" json:"role"`
	Name       string   `json:"name"`
	Phone      string   `json:"phone"`
	IsActive   bool     `gorm:"default:true" json:"is_active"`
	FirstLogin bool     `gorm:"default:true" json:"first_login"` // 是否首次登录
}
