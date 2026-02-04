package controllers

import "github.com/gin-gonic/gin"

// 认证相关
func Login(c *gin.Context)          {}
func Register(c *gin.Context)       {} // 注册 (初始化用)
func UpdatePassword(c *gin.Context) {}
func Logout(c *gin.Context)         {}

// 用户/学生/老师相关
func GetProfile(c *gin.Context) {}

// 物品相关
func CreateItem(c *gin.Context)   {}
func UpdateItem(c *gin.Context)   {}
func DeleteItem(c *gin.Context)   {} // 取消或删除
func GetItem(c *gin.Context)      {}
func GetItemsList(c *gin.Context) {}
func GetMyItems(c *gin.Context)   {}

// 认领相关
func CreateClaim(c *gin.Context) {}
func GetClaims(c *gin.Context)   {}

// 管理员相关
func AdminAuditItem(c *gin.Context)        {}
func AdminGetItems(c *gin.Context)         {}
func AdminGetUserList(c *gin.Context)      {}
func AdminCreateUser(c *gin.Context)       {} // 添加失物招领管理员
func AdminUpdateUserStatus(c *gin.Context) {} // 禁用账号

// Upload
func UploadImage(c *gin.Context) {}
