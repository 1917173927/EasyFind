package ws

import "time"

type UpdateScope string

const (
	ScopeDialog        UpdateScope = "dialog"
	ScopeActivity      UpdateScope = "activity"
	ScopeClaimProgress UpdateScope = "claim_progress"
	ScopeAnnouncement  UpdateScope = "announcement"
)

type UpdateEvent struct {
	Type        string      `json:"type"`
	Scope       UpdateScope `json:"scope"`
	HasNew      bool        `json:"has_new"`
	UnreadCount *int64      `json:"unread_count,omitempty"`
	RefID       uint        `json:"ref_id,omitempty"`
	At          time.Time   `json:"at"`
}

func PushUpdate(userID uint, scope UpdateScope, unreadCount *int64, refID uint) error {
	event := UpdateEvent{
		Type:        "update",
		Scope:       scope,
		HasNew:      true,
		UnreadCount: unreadCount,
		RefID:       refID,
		At:          time.Now(),
	}
	return Manager.SafeWriteJSON(userID, event)
}

func BroadcastUpdate(scope UpdateScope, refID uint) {
	event := UpdateEvent{
		Type:   "update",
		Scope:  scope,
		HasNew: true,
		RefID:  refID,
		At:     time.Now(),
	}
	Manager.BroadcastJSON(event)
}
