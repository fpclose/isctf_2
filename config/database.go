package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDatabase 连接数据库
func ConnectDatabase() {
	// 获取DSN
	dsn := GetDSN()

	// GORM 配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置日志级别
		// 禁用外键约束（可选）
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	// 连接数据库
	database, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 获取底层 SQL 数据库
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("获取数据库实例失败:", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(0)         // 连接最大生命周期

	DB = database
	fmt.Println("数据库连接成功")
}

// AutoMigrate 自动迁移数据表
func AutoMigrate(models ...interface{}) {
	err := DB.AutoMigrate(models...)
	if err != nil {
		log.Fatal("数据表迁移失败:", err)
	}
	fmt.Println("数据表迁移成功")
}
