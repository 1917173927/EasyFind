package models

import (
	"time"

	"gorm.io/gorm"
)

type ItemType string
type ItemStatus string

const (
	TypeLost  ItemType = "lost"
	TypeFound ItemType = "found"

	StatusPending   ItemStatus = "pending"
	StatusApproved  ItemStatus = "approved"
	StatusMatched   ItemStatus = "matched"
	StatusClaimed   ItemStatus = "claimed"
	StatusRejected  ItemStatus = "rejected"
	StatusCancelled ItemStatus = "cancelled"
	StatusArchived  ItemStatus = "archived"
)

type Item struct {
	gorm.Model
	Title        string     `gorm:"not null" json:"title"`                                                                                                       // 标题
	Type         ItemType   `gorm:"type:enum('lost', 'found');not null" json:"type"`                                                                             // 类型 (lost: 失物, found: 招领)
	Category     string     `json:"category"`                                                                                                                    // 分类
	Location     string     `json:"location"`                                                                                                                    // 地点
	Time         time.Time  `json:"time"`                                                                                                                        // 时间
	Description  string     `json:"description"`                                                                                                                 // 描述
	Photos       string     `json:"photos"`                                                                                                                      // 照片 (JSON或逗号分隔的URL)
	ContactName  string     `json:"contact_name"`                                                                                                                // 联系人姓名
	ContactPhone string     `json:"contact_phone"`                                                                                                               // 联系人电话
	IsBounty     bool       `json:"is_bounty"`                                                                                                                   // 是否有悬赏(保留字段)
	Status       ItemStatus `gorm:"type:enum('pending', 'approved', 'matched', 'claimed', 'rejected', 'cancelled', 'archived');default:'pending'" json:"status"` // 状态

	PublisherID uint    `json:"publisher_id"`
	Publisher   Account `gorm:"foreignKey:PublisherID" json:"publisher"` // 发布者信息

	ReviewerID *uint    `json:"reviewer_id"`
	Reviewer   *Account `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"` // 审核者信息

	RejectReason string `json:"reject_reason"` // 审核拒绝原因
}

type Claim struct {
	gorm.Model
	ItemID     uint    `gorm:"not null" json:"item_id"`
	Item       Item    `gorm:"foreignKey:ItemID" json:"item"`
	ClaimantID uint    `gorm:"not null" json:"claimant_id"`
	Claimant   Account `gorm:"foreignKey:ClaimantID" json:"claimant"`
	Status     string  `gorm:"default:'pending'" json:"status"` // 状态: pending/approved/rejected
	Proof      string  `json:"proof"`                           // 证明材料 (如图片URL)
}
