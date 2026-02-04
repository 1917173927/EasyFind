package controllers

import "github.com/gin-gonic/gin"

// 审核操作 (通过/驳回)
func AdminAuditItem(c *gin.Context) {}

// 待审核列表
func AdminGetItems(c *gin.Context) {}

// 用户列表
func AdminGetUserList(c *gin.Context) {}

// 添加失物招领管理员
func AdminCreateUser(c *gin.Context) {}

// 禁用账号
func AdminUpdateUserStatus(c *gin.Context) {}
