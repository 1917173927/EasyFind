package database

import (
	"fmt"
	"log"

	"easyfind/internal/config"
	"easyfind/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMySQL() {
	cfg := config.AppConfigData.Database.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 可选: 自动迁移模型
	err = DB.AutoMigrate(&models.Account{}, &models.Item{}, &models.Claim{}, &models.LostCategory{}, &models.Image{})
	if err != nil {
		log.Printf("auto migrate failed: %v", err)
	}

	// 初始化分类数据 (如果表为空)
	var count int64
	DB.Model(&models.LostCategory{}).Count(&count)
	if count == 0 {
		categories := []models.LostCategory{
			{CategoryName: "电子产品"},
			{CategoryName: "生活用品"},
			{CategoryName: "书籍文具"},
			{CategoryName: "证件卡片"},
			{CategoryName: "衣物服饰"},
			{CategoryName: "其他"},
		}
		DB.Create(&categories)
		log.Println("Initialized default categories")
	}
}
