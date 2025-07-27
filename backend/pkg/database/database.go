package database

import (
	"github.com/3d-printing-farm/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移数据库表
	err = db.AutoMigrate(
		&models.User{},
		&models.Material{},
		&models.File{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		return nil, err
	}

	// 初始化基础数据
	initializeData(db)

	return db, nil
}

func initializeData(db *gorm.DB) {
	// 初始化材料数据
	materials := []models.Material{
		{Name: "PLA", Description: "环保易用，适合初学者", PricePerGram: 0.8, Properties: `{"temperature": "190-220°C", "flexibility": "低"}`},
		{Name: "ABS", Description: "耐用性强，适合功能件", PricePerGram: 1.0, Properties: `{"temperature": "220-260°C", "flexibility": "中"}`},
		{Name: "PETG", Description: "透明度好，化学稳定", PricePerGram: 1.2, Properties: `{"temperature": "220-250°C", "flexibility": "中"}`},
		{Name: "TPU", Description: "柔性材料，适合弹性件", PricePerGram: 2.0, Properties: `{"temperature": "210-230°C", "flexibility": "高"}`},
		{Name: "WOOD", Description: "木质纤维复合材料", PricePerGram: 1.5, Properties: `{"temperature": "190-220°C", "flexibility": "低"}`},
	}

	for _, material := range materials {
		var existing models.Material
		if err := db.Where("name = ?", material.Name).First(&existing).Error; err != nil {
			db.Create(&material)
		}
	}
}