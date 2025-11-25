package config

import (
	"fmt"
	"os"
)

// Config 全局配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Mode string // debug, release, test
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Charset  string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string
	ExpireTime int // 小时
}

var AppConfig *Config

// InitConfig 初始化配置
func InitConfig() {
	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8081"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			Username: getEnv("DB_USERNAME", "root"),
			Password: getEnv("DB_PASSWORD", "123456"),
			Database: getEnv("DB_DATABASE", "dalictf"),
			Charset:  getEnv("DB_CHARSET", "utf8mb4"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "isctf-secret-key-2024"),
			ExpireTime: 24, // 24小时
		},
	}

	fmt.Println("配置加载成功")
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDSN 获取数据库连接字符串
func GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		AppConfig.Database.Username,
		AppConfig.Database.Password,
		AppConfig.Database.Host,
		AppConfig.Database.Port,
		AppConfig.Database.Database,
		AppConfig.Database.Charset,
	)
}
