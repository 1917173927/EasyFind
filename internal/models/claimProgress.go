package models

import "time"

const (
	ClaimProgressTypePending   = "claim_pending"
	ClaimProgressTypeApproved  = "claim_approved"
	ClaimProgressTypeRejected  = "claim_rejected"
	ClaimProgressTypeCompleted = "claim_completed"
)

// ClaimProgress 招领进度视图模型
type ClaimProgress struct {
	Type         string    `json:"type"`                    // 进度类型: claim_pending|claim_approved|claim_rejected|claim_completed
	Title        string    `json:"title"`                   // 主标题
	Hint         string    `json:"hint,omitempty"`          // 辅助提示
	ActionText   string    `json:"action_text,omitempty"`   // 按钮文案
	Time         time.Time `json:"time"`                    // 动态时间（倒序）
	ClaimID      uint      `json:"claim_id"`                // 认领记录ID（不是用户ID）
	ItemID       uint      `json:"item_id"`                 // 物品ID
	PeerUserID   uint      `json:"peer_user_id"`            // 对方用户ID（物品发布者ID，可用于发起会话）
	ItemName     string    `json:"item_name"`               // 物品名称
	LossTime     time.Time `json:"loss_time"`               // 物品时间
	Location     string    `json:"location"`                // 地点
	Img          string    `json:"img"`                     // 缩略图
	ClaimStatus  string    `json:"claim_status"`            // 认领状态: pending|approved|rejected|archived
	RejectReason string    `json:"reject_reason,omitempty"` // 驳回原因
}
