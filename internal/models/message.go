package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	SenderID   uint    `gorm:"index;not null" json:"sender_id"`
	Sender     Account `gorm:"foreignKey:SenderID" json:"sender"`
	ReceiverID uint    `gorm:"index;not null" json:"receiver_id"`
	Receiver   Account `gorm:"foreignKey:ReceiverID" json:"receiver"`
	Content    string  `gorm:"not null" json:"content"`
	Type       int     `gorm:"default:1" json:"type"` //1文本 2图片
	IsRead     bool    `gorm:"default:false" json:"is_read"`
	ItemID     uint    `gorm:"index;default:0" json:"item_id"`
}

type ChatList struct {
	TargetID    uint      `json:"target_id"`    // 对方的用户ID
	TargetName  string    `json:"target_name"`  // 对方昵称
	Avatar      string    `json:"avatar"`       // 对方头像
	LastMsg     string    `json:"last_msg"`     // 最后一条消息内容
	LastTime    time.Time `json:"last_time"`    // 最后一条消息时间
	UnreadCount int64     `json:"unread_count"` // 红点数量
}
