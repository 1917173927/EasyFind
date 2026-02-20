package services

import (
	"easyfind/internal/models"
	"easyfind/internal/ws"
	"easyfind/pkg/database"
	"sort"

	"go.uber.org/zap"
)

// SendMessage 发送消息
func SendMessage(senderID, receiverID uint, content string, msgType int, itemID uint) error {
	msg := models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Type:       msgType,
		ItemID:     itemID,
		IsRead:     false,
	}

	if err := database.DB.Create(&msg).Error; err != nil {
		return err
	}

	var fullMsg models.Message
	database.DB.Preload("Sender").Preload("Receiver").First(&fullMsg, msg.ID)

	err := ws.Manager.SafeWriteJSON(receiverID, fullMsg)
	if err != nil {
		zap.L().Warn("Failed to push message via WebSocket",
			zap.Uint("receiver_id", receiverID),
			zap.Error(err))
	}

	var unreadCount int64
	database.DB.Model(&models.Message{}).Where("receiver_id = ? AND is_read = ?", receiverID, false).Count(&unreadCount)
	_ = ws.PushUpdate(receiverID, ws.ScopeDialog, &unreadCount, msg.ID)

	return nil
}

// GetHistoryMessages 获取历史消息
func GetHistoryMessages(userID, targetID uint, cursor uint, limit int) ([]models.Message, error) {
	var messages []models.Message

	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	query := database.DB.Model(&models.Message{}).
		Preload("Sender").
		Preload("Receiver").
		Where(
			database.DB.Where("sender_id = ? AND receiver_id = ?", userID, targetID).
				Or("sender_id = ? AND receiver_id = ?", targetID, userID),
		)

	if cursor > 0 {
		query = query.Where("id < ?", cursor)
	}

	err := query.Order("id desc").Limit(limit).Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}

// MarkMessagesAsRead 标记消息为已读
func MarkMessagesAsRead(userID, targetID uint) error {
	err := database.DB.Model(&models.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?", targetID, userID, false).
		Update("is_read", true).Error
	if err != nil {
		return err
	}

	var unreadCount int64
	database.DB.Model(&models.Message{}).Where("receiver_id = ? AND is_read = ?", userID, false).Count(&unreadCount)
	_ = ws.PushUpdate(userID, ws.ScopeDialog, &unreadCount, 0)

	return nil
}

// GetChatList 获取会话列表
func GetChatList(userID uint) ([]models.ChatList, error) {
	var targetIDs []uint

	err := database.DB.Raw(`
		SELECT DISTINCT target_id FROM (
			SELECT CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END AS target_id, created_at
			FROM messages 
			WHERE (sender_id = ? OR receiver_id = ?) AND deleted_at IS NULL
			ORDER BY created_at DESC
			LIMIT 200 
		) AS t LIMIT 50`, userID, userID, userID).Scan(&targetIDs).Error

	if err != nil {
		return nil, err
	}

	if len(targetIDs) == 0 {
		return []models.ChatList{}, nil
	}

	var users []models.Account
	if err := database.DB.Where("id IN ?", targetIDs).Find(&users).Error; err != nil {
		return nil, err
	}
	userMap := make(map[uint]models.Account)
	for _, u := range users {
		userMap[u.ID] = u
	}

	type UnreadResult struct {
		SenderID uint
		Count    int64
	}
	var unreads []UnreadResult
	database.DB.Model(&models.Message{}).
		Select("sender_id, count(*) as count").
		Where("receiver_id = ? AND is_read = ? AND sender_id IN ?", userID, false, targetIDs).
		Group("sender_id").
		Scan(&unreads)

	unreadMap := make(map[uint]int64)
	for _, u := range unreads {
		unreadMap[u.SenderID] = u.Count
	}

	var results []models.ChatList
	for _, targetID := range targetIDs {
		targetUser, exists := userMap[targetID]
		if !exists {
			continue
		}

		item := models.ChatList{
			TargetID:    targetUser.ID,
			TargetName:  targetUser.Username,
			Avatar:      targetUser.Avatar,
			UnreadCount: unreadMap[targetID],
		}

		var lastMsg models.Message

		err := database.DB.Where(
			"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, targetID, targetID, userID,
		).Order("created_at desc").First(&lastMsg).Error

		if err == nil {
			item.LastTime = lastMsg.CreatedAt
			if lastMsg.Type == 2 {
				item.LastMsg = "[图片]"
			} else {
				// 简单的截断逻辑
				runes := []rune(lastMsg.Content)
				if len(runes) > 30 {
					item.LastMsg = string(runes[:30]) + "..."
				} else {
					item.LastMsg = lastMsg.Content
				}
			}
		}

		results = append(results, item)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].LastTime.After(results[j].LastTime)
	})

	return results, nil
}
