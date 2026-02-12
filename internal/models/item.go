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

	StatusPending   ItemStatus = "pending"   // 待审核
	StatusApproved  ItemStatus = "approved"  // 已通过
	StatusMatched   ItemStatus = "matched"   // 已匹配（找到对应的人来对话了）
	StatusClaimed   ItemStatus = "claimed"   // 已认领（确认是对的人）
	StatusRejected  ItemStatus = "rejected"  // 已驳回
	StatusCancelled ItemStatus = "cancelled" // 已取消
	StatusArchived  ItemStatus = "archived"  // 已归档（认领成功或者超时未认领）
)

type Item struct {
	gorm.Model
	Title        string     `gorm:"not null" json:"title"` // 标题
	Type         ItemType   `gorm:"type:enum('lost', 'found');not null" json:"type"`
	Campus       string     `json:"campus"`      // 类型 (lost: 失物, found: 招领)
	Category     string     `json:"category"`    // 分类
	Location     string     `json:"location"`    // 地点
	Time         time.Time  `json:"time"`        // 时间
	Description  string     `json:"description"` // 描述
	Img1         string     `json:"img1"`        // 照片 (JSON或逗号分隔的URL)
	Img2         string     `json:"img2"`        // 照片 (JSON或逗号分隔的URL)
	Img3         string     `json:"img3"`        // 照片 (JSON或逗号分隔的URL)
	Img4         string     `json:"img4"`
	ContactName  string     `json:"contact_name"`                                                                                                                // 联系人姓名
	ContactPhone string     `json:"contact_phone"`                                                                                                               // 联系人电话
	Bounty       float64    `gorm:"default:0" json:"bounty"`                                                                                                     // 悬赏金额 (0表示无悬赏)
	Status       ItemStatus `gorm:"type:enum('pending', 'approved', 'matched', 'claimed', 'rejected', 'cancelled', 'archived');default:'pending'" json:"status"` // 状态

	PublisherID uint    `json:"publisher_id"`
	Publisher   Account `gorm:"foreignKey:PublisherID" json:"publisher"` // 发布者信息

	ReviewerID *uint    `json:"reviewer_id"`
	Reviewer   *Account `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"` // 审核者信息

	RejectReason string `json:"reject_reason"` // 审核拒绝原因

	ProcessMethod string `json:"process_method"` // 处理方式
}

type Claim struct {
	gorm.Model
	ItemID     uint    `gorm:"not null" json:"item_id"`
	Item       Item    `gorm:"foreignKey:ItemID" json:"item"`
	ClaimantID uint    `gorm:"not null" json:"claimant_id"`
	Claimant   Account `gorm:"foreignKey:ClaimantID" json:"claimant"`
	Status     string  `gorm:"default:'pending'" json:"status"` // 状态: pending/approved/rejected/archived
	Proof      string  `json:"proof"`                           // 证明材料 (如图片URL)
}

type LostCategory struct {
	ID           int    `json:"id"`
	CategoryName string `json:"kind_name"`
}
