package models

import "gorm.io/gorm"

// Announcement 系统公告
type Announcement struct {
	gorm.Model
	Title     string `gorm:"not null" json:"title"`
	Content   string `gorm:"type:text;not null" json:"content"`
	Type      string `gorm:"default:'global'" json:"type"`      // global: 全局公告, region: 区域公告
	Publisher string `json:"publisher"`                         // 发布者姓名
	IsTop     bool   `json:"is_top"`                            // 是否置顶
	Status    string `gorm:"default:'published'" json:"status"` // published, pending, rejected
	Region    string `json:"region"`                            // 校区 (仅Type=region有效)
}

// Feedback 用户反馈/投诉
type Feedback struct {
	gorm.Model
	UserID  uint    `json:"user_id"`
	User    Account `gorm:"foreignKey:UserID" json:"user"`
	Type    string  `json:"type"` // bug, complaint, suggestion
	Content string  `gorm:"type:text" json:"content"`
	Contact string  `json:"contact"`                         // 联系方式
	Status  string  `gorm:"default:'pending'" json:"status"` // pending, resolved, ignored
	Reply   string  `gorm:"type:text" json:"reply"`          // 管理员回复
	Handler string  `json:"handler"`                         // 处理人
}

// SystemStats 系统统计数据
type SystemStats struct {
	TotalUsers  int64 `json:"total_users"`
	ActiveUsers int64 `json:"active_users"` // 30天内活跃
	TotalItems  int64 `json:"total_items"`
	SolvedItems int64 `json:"solved_items"` // status = claimed/matched
	TotalClaims int64 `json:"total_claims"`
	TodayItems  int64 `json:"today_items"`
}
