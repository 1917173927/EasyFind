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
	Username   string   `gorm:"uniqueIndex;not null" json:"username"`                                                        // 学号或工号
	Password   string   `gorm:"not null" json:"-"`                                                                           // 密码
	Role       UserRole `gorm:"type:enum('student_teacher', 'lf_admin', 'sys_admin');default:'student_teacher'" json:"role"` // 角色
	Name       string   `json:"name"`                                                                                        // 姓名
	Phone      string   `json:"phone"`                                                                                       // 电话号码
	IsActive   bool     `gorm:"default:true" json:"is_active"`                                                               // 账号是否启用
	FirstLogin bool     `gorm:"default:true" json:"first_login"`                                                             // 是否首次登录
}
