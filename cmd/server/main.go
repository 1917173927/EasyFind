package main

import (
	"log"

	"easyfind/internal/api"
	"easyfind/internal/config"
	"easyfind/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 数据库初始化
	database.InitMySQL()
	database.InitRedis()

	// 路由初始化
	r := gin.Default()
	api.Init(r)

	// 服务启动
	port := config.AppConfigData.Server.Port
	log.Printf("Starting server on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
