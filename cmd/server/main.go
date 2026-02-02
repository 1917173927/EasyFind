package main

import (
	"log"

	"easyfind/internal/api"
	"easyfind/internal/config"
	"easyfind/pkg/database"
)

func main() {
	// 1. Init Config
	config.InitConfig()

	// 2. Init Database
	database.InitMySQL()
	database.InitRedis()

	// 3. Setup Router
	r := api.SetupRouter()

	// 4. Run Server
	port := config.AppConfigData.Server.Port
	log.Printf("Starting server on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
