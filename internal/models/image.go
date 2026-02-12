package models

import "gorm.io/gorm"

// Image 图片存储模型 (去重)
type Image struct {
	gorm.Model
	FileName string `gorm:"not null" json:"file_name"`                         // 原始文件名
	Hash     string `gorm:"type:varchar(64);uniqueIndex;not null" json:"hash"` // 文件哈希 (MD5), 用于去重
	URL      string `gorm:"not null" json:"url"`                               // 访问 URL
	Path     string `gorm:"not null" json:"-"`                                 // 本地存储路径
	MimeType string `json:"mime_type"`                                         // 文件类型
	Size     int64  `json:"size"`                                              // 文件大小
}
