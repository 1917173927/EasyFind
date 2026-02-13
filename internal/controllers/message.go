package controllers

import (
	"easyfind/internal/apiErr"
	"easyfind/internal/services"
	"easyfind/pkg/response"

	"github.com/gin-gonic/gin"
)

type SendMessageRequest struct {
	ReceiverID uint   `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Type       int    `json:"type" binding:"omitempty,oneof=1 2"` // 1:文本, 2:图片
	ItemID     uint   `json:"item_id"`
}

// SendMessage 发送消息
func SendMessage(c *gin.Context) {
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	uid, exists := c.Get("userID")
	if !exists {
		apiErr.HandleBizError(c, response.ErrTokenInvalid, "未获取到用户信息")
		return
	}
	userID, ok := uid.(uint)
	if !ok {
		apiErr.HandleSysError(c, response.CodeError, nil) // 类型断言失败视为系统错误
		return
	}

	if req.Type == 0 {
		req.Type = 1 // 默认文本
	}

	if err := services.SendMessage(userID, req.ReceiverID, req.Content, req.Type, req.ItemID); err != nil {
		apiErr.HandleSysError(c, response.ErrMessageSendFail, err)
		return
	}

	response.Success(c, nil)
}

type GetHistoryMessagesRequest struct {
	TargetID uint `form:"target_id" binding:"required"`
	Cursor   uint `form:"cursor"`
	Limit    int  `form:"limit"`
}

// GetHistoryMessages 获取历史消息
func GetHistoryMessages(c *gin.Context) {
	var req GetHistoryMessagesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	uid, exists := c.Get("userID")
	if !exists {
		apiErr.HandleBizError(c, response.ErrTokenInvalid, "未获取到用户信息")
		return
	}
	userID, ok := uid.(uint)
	if !ok {
		apiErr.HandleSysError(c, response.CodeError, nil)
		return
	}

	messages, err := services.GetHistoryMessages(userID, req.TargetID, req.Cursor, req.Limit)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrMessageGetFail, err)
		return
	}

	response.Success(c, messages)
}

// GetChatList 获取会话列表
func GetChatList(c *gin.Context) {
	uid, exists := c.Get("userID")
	if !exists {
		apiErr.HandleBizError(c, response.ErrTokenInvalid, "未获取到用户信息")
		return
	}
	userID, ok := uid.(uint)
	if !ok {
		apiErr.HandleSysError(c, response.CodeError, nil)
		return
	}

	chatList, err := services.GetChatList(userID)
	if err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	response.Success(c, chatList)
}

type MarkReadRequest struct {
	TargetID uint `json:"target_id" binding:"required"`
}

// MarkMessagesAsRead 标记消息为已读
func MarkMessagesAsRead(c *gin.Context) {
	var req MarkReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr.HandleValidatorError(c, err)
		return
	}

	uid, exists := c.Get("userID")
	if !exists {
		apiErr.HandleBizError(c, response.ErrTokenInvalid, "未获取到用户信息")
		return
	}
	userID, ok := uid.(uint)
	if !ok {
		apiErr.HandleSysError(c, response.CodeError, nil)
		return
	}

	if err := services.MarkMessagesAsRead(userID, req.TargetID); err != nil {
		apiErr.HandleSysError(c, response.ErrDBQueryFail, err)
		return
	}

	response.Success(c, nil)
}
