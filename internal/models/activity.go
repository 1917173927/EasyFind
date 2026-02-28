package models

import "time"

// PostActivity 帖子动态视图模型
type PostActivity struct {
	Type         string    `json:"type"`                    // status_update (状态更新), claim_received (收到认领), rejected (驳回)
	Title        string    `json:"title"`                   // 标题
	Time         time.Time `json:"time"`                    // 发生时间
	ItemID       uint      `json:"item_id"`                 // 物品ID
	PeerUserID   uint      `json:"peer_user_id"`            // 对方用户ID（无对方时为0）
	ItemName     string    `json:"item_name"`               // 物品名称 (Title)
	LossTime     time.Time `json:"loss_time"`               // 物品关联的时间(丢失/拾取时间)
	Location     string    `json:"location"`                // 地点
	Img          string    `json:"img"`                     // 缩略图
	Status       string    `json:"item_status"`             // 物品当前状态 (pending, approved, etc.)
	RejectReason string    `json:"reject_reason,omitempty"` // 驳回原因
	ClaimID      uint      `json:"claim_id,omitempty"`      // 认领ID
	ClaimantName string    `json:"claimant_name,omitempty"` // 认领人名称
}
