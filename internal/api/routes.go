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
	// 设置上传限制
	r.MaxMultipartMemory = 64 << 20 // 64MB
	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// Swagger 文档路由访问 http://localhost:8080/swagger/index.html 即可查看文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	

	// 全局前缀
	api := r.Group("/api/v1")
	{
		// 公开接口 (无需 JWT)
		api.POST("/login", controllers.Login)
		api.POST("/register", controllers.Register) // 临时注册接口

		// 需要 JWT 认证的接口
		auth := api.Group("")
		auth.Use(middleware.JWT())
		{
			// 用户相关
			auth.POST("/auth/logout", controllers.Logout)
			auth.PUT("/auth/password", controllers.UpdatePassword)
			auth.GET("/user/profile", controllers.GetProfile)

			// 物品发布与管理 (学生/老师)
			auth.POST("/items", controllers.CreateItem)       // 发布物品
			auth.PUT("/items/:id", controllers.UpdateItem)    // 修改发布
			auth.DELETE("/items/:id", controllers.DeleteItem) // 取消/删除发布
			auth.GET("/items", controllers.GetItemsList)      // 查询物品列表 (筛选)
			auth.GET("/items/:id", controllers.GetItem)       // 查看详情
			auth.GET("/my/items", controllers.GetMyItems)     // 我的发布

			// 认领与沟通
			auth.POST("/claims", controllers.CreateClaim) // 提交认领
			auth.GET("/claims", controllers.GetClaims)    // 查看认领进度/列表

			// 图片上传
			auth.POST("/upload", controllers.UploadImage)

			// 管理员接口 (需要额外的权限检查中间件，此处简化放在一起)
			admin := auth.Group("/admin")
			{
				// 失物招领管理员
				admin.GET("/audit/items", controllers.AdminGetItems)      // 待审核列表
				admin.PUT("/audit/items/:id", controllers.AdminAuditItem) // 审核操作 (通过/驳回)

				// 系统管理员
				admin.GET("/users", controllers.AdminGetUserList)                 // 用户列表
				admin.POST("/users", controllers.AdminCreateUser)                 // 新增管理员
				admin.PUT("/users/:id/status", controllers.AdminUpdateUserStatus) // 禁用/解禁账号
			}
		}
	}
}
