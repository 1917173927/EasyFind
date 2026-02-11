package api

import (
	"easyfind/internal/controllers"
	"easyfind/internal/middleware"

	_ "easyfind/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init(r *gin.Engine) {
	r.MaxMultipartMemory = 64 << 20
	r.Static("/uploads", "./uploads")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		// 公开接口 (Public)
		api.POST("/login", controllers.Login)
		api.POST("/register", controllers.Register) // 仅用于测试

		// 物品浏览
		api.GET("/items/all", controllers.GetAllItem) // 获取所有物品 (不分 Lost/Found)
		api.GET("/items", controllers.GetRecord)      // 获取物品 (区分 Lost/Found)
		api.GET("/items/:id", controllers.GetRecordById)
		api.GET("/kinds", controllers.GetCategoryList)

		// 需要认证的通用接口 (User)
		auth := api.Group("")
		auth.Use(middleware.JWT())
		{
			// 个人信息
			auth.POST("/logout", controllers.Logout)
			auth.PUT("/password", controllers.UpdatePassword)
			auth.GET("/profile", controllers.GetProfile)
			auth.PUT("/profile", controllers.UpdateProfile)
			auth.DELETE("/account", controllers.DeleteAccount)

			// 我的物品管理
			auth.POST("/items", controllers.CreateRecord)
			auth.PUT("/items/:id", controllers.UpdateRecord)
			auth.DELETE("/items/:id", controllers.DeleteRecord)
			auth.PUT("/items/:id/cancel", controllers.CancelRecord)
			auth.GET("/my/items", controllers.GetAllMyRecord)

			// 我的认领管理
			auth.POST("/claims", controllers.CreatClaim)
			auth.GET("/my/claims", controllers.GetMyClaim)
			auth.GET("/claims/:id", controllers.GetClaimByID)
			auth.PUT("/claims/:id/confirm", controllers.ConfirmClaim)

			// 图片上传
			// auth.POST("/upload", controllers.UploadImage)
		}

		//  管理员接口 (Admin)
		admin := api.Group("/admin")
		admin.Use(middleware.JWT(), middleware.AuthAdmin())
		{
			// 物品管理
			admin.GET("/items", controllers.GetRecordByAdmin)
			admin.GET("/items/pending", controllers.GetPendingRecordByAdmin)
			admin.PUT("/items/:id/approve", controllers.ApproveRecord)
			admin.PUT("/items/:id/reject", controllers.RejectRecord)
			admin.PUT("/items/:id/archive", controllers.ArchiveRecord)

			// 认领管理
			admin.GET("/claims/pending", controllers.GetPendingClaimByAdmin)
			admin.PUT("/claims/:id/approve", controllers.ApproveClaim)
			admin.PUT("/claims/:id/reject", controllers.RejectClaim)
		}

		// 超级管理员接口 (System Admin)
		sysAdmin := api.Group("/sysadmin")
		sysAdmin.Use(middleware.JWT(), middleware.AuthSysAdmin())
		{
			// 用户管理
			// sysAdmin.GET("/users", controllers.AdminGetUserList)
			// sysAdmin.POST("/users", controllers.AdminCreateUser)
			// sysAdmin.PUT("/users/:id/status", controllers.AdminUpdateUserStatus)
		}
	}
}
