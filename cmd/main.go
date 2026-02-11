package main

import (
	"log"

	"easyfind/internal/api"
	"easyfind/internal/config"
	"easyfind/pkg/database"

	"github.com/gin-gonic/gin"
)

// @title           EasyFind API
// @version         1.0
// @description     This is the API server for EasyFind application.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

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
