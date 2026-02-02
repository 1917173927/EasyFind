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
	Title        string     `gorm:"not null" json:"title"`
	Type         ItemType   `gorm:"type:enum('lost', 'found');not null" json:"type"`
	Category     string     `json:"category"`
	Location     string     `json:"location"`
	Time         time.Time  `json:"time"`
	Description  string     `json:"description"`
	Photos       string     `json:"photos"` // 照片
	ContactName  string     `json:"contact_name"`
	ContactPhone string     `json:"contact_phone"`
	IsBounty     bool       `json:"is_bounty"`
	Status       ItemStatus `gorm:"type:enum('pending', 'approved', 'matched', 'claimed', 'rejected', 'cancelled', 'archived');default:'pending'" json:"status"`

	PublisherID uint `json:"publisher_id"`
	Publisher   User `gorm:"foreignKey:PublisherID" json:"publisher"`

	ReviewerID *uint `json:"reviewer_id"`
	Reviewer   *User `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`

	RejectReason string `json:"reject_reason"`
}

type Claim struct {
	gorm.Model
	ItemID     uint   `gorm:"not null" json:"item_id"`
	Item       Item   `gorm:"foreignKey:ItemID" json:"item"`
	ClaimantID uint   `gorm:"not null" json:"claimant_id"`
	Claimant   User   `gorm:"foreignKey:ClaimantID" json:"claimant"`
	Status     string `gorm:"default:'pending'" json:"status"` // 状态: pending/approved/rejected
	Proof      string `json:"proof"`
}
