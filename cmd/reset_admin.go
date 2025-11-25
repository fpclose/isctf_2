package cmd

import (
	"fmt"
	"isctf/config"
	"isctf/models"
	"log"
)

// ResetAdminPassword 重置管理员密码
func ResetAdminPassword() {
	var admin models.User
	
	// 查找 root_admin 用户
	if err := config.DB.Where("username = ?", "root_admin").First(&admin).Error; err != nil {
		log.Printf("查找管理员失败: %v", err)
		return
	}

	// 更新密码
	if err := admin.SetPassword("fpclose_SfTian_i5ctf"); err != nil {
		log.Printf("设置密码失败: %v", err)
		return
	}

	// 保存到数据库
	if err := config.DB.Model(&admin).Update("password", admin.Password).Error; err != nil {
		log.Printf("更新密码失败: %v", err)
		return
	}

	fmt.Println("✓ 管理员密码重置成功")
	fmt.Printf("  用户名: %s\n", admin.Username)
	fmt.Printf("  新密码: %s\n", "fpclose_SfTian_i5ctf")
	fmt.Printf("  邮箱: %s\n", admin.Email)
	fmt.Printf("  角色: %s\n", admin.Role)
}

// DeleteAndRecreateAdmin 删除并重新创建管理员
func DeleteAndRecreateAdmin() {
	// 删除现有的 root_admin
	if err := config.DB.Where("username = ?", "root_admin").Delete(&models.User{}).Error; err != nil {
		log.Printf("删除现有管理员失败: %v", err)
		return
	}

	fmt.Println("✓ 已删除现有的管理员账户")

	// 创建新的管理员
	admin := &models.User{
		Username:      "root_admin",
		Email:         "admin@isctf.com",
		Role:          "super_admin",
		Track:         "social",
		EmailVerified: true,
		VerifyStatus:  "approved",
		Status:        "active",
	}

	// 设置密码
	if err := admin.SetPassword("fpclose_SfTian_i5ctf"); err != nil {
		log.Printf("设置管理员密码失败: %v", err)
		return
	}

	// 保存到数据库
	if err := config.DB.Create(admin).Error; err != nil {
		log.Printf("创建管理员账户失败: %v", err)
		return
	}

	fmt.Println("✓ 新管理员账户创建成功")
	fmt.Printf("  用户名: %s\n", admin.Username)
	fmt.Printf("  密码: %s\n", "fpclose_SfTian_i5ctf")
	fmt.Printf("  邮箱: %s\n", admin.Email)
	fmt.Printf("  角色: %s\n", admin.Role)
}
