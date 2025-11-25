package cmd

import (
	"fmt"
	"isctf/config"
	"isctf/models"
	"log"

	"gorm.io/gorm"
)

// InitDefaultAdmin 初始化默认管理员账户
func InitDefaultAdmin() {
	// 检查是否已存在 root_admin
	var existingAdmin models.User
	err := config.DB.Where("username = ?", "root_admin").First(&existingAdmin).Error
	if err == nil {
		fmt.Println("✓ 默认管理员账户已存在")
		return
	}

	if err != gorm.ErrRecordNotFound {
		log.Printf("检查管理员账户时出错: %v", err)
		return
	}

	// 创建默认管理员
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

	fmt.Println("✓ 默认管理员账户创建成功")
	fmt.Printf("  用户名: %s\n", admin.Username)
	fmt.Printf("  密码: %s\n", "fpclose_SfTian_i5ctf")
	fmt.Printf("  邮箱: %s\n", admin.Email)
	fmt.Printf("  角色: %s\n", admin.Role)
}

// InitTestData 初始化测试数据
func InitTestData() {
	// 检查是否已有测试学校
	var school models.School
	err := config.DB.Where("school_name = ?", "测试大学").First(&school).Error
	if err == nil {
		fmt.Println("✓ 测试数据已存在")
		return
	}

	// 创建测试学校
	school = models.School{
		SchoolName: "测试大学",
		Status:     "active",
	}
	if err := config.DB.Create(&school).Error; err != nil {
		log.Printf("创建测试学校失败: %v", err)
		return
	}

	// 创建测试管理员（学校管理员）
	schoolAdmin := &models.User{
		Username:      "school_admin",
		Email:         "admin@test.edu.cn",
		Role:          "school_admin",
		Track:         "school",
		SchoolID:      &school.ID,
		SchoolName:    &school.SchoolName,
		EmailVerified: true,
		VerifyStatus:  "approved",
		Status:        "active",
	}

	if err := schoolAdmin.SetPassword("test123456!"); err != nil {
		log.Printf("设置学校管理员密码失败: %v", err)
		return
	}

	if err := config.DB.Create(schoolAdmin).Error; err != nil {
		log.Printf("创建学校管理员失败: %v", err)
		return
	}

	// 创建测试用户
	testUser := &models.User{
		Username:      "test_user",
		Email:         "user@test.com",
		Role:          "user",
		Track:         "social",
		EmailVerified: true,
		VerifyStatus:  "approved",
		Status:        "active",
	}

	if err := testUser.SetPassword("test123456!"); err != nil {
		log.Printf("设置测试用户密码失败: %v", err)
		return
	}

	if err := config.DB.Create(testUser).Error; err != nil {
		log.Printf("创建测试用户失败: %v", err)
		return
	}

	fmt.Println("✓ 测试数据创建成功")
	fmt.Printf("  学校: %s (ID: %d)\n", school.SchoolName, school.ID)
	fmt.Printf("  学校管理员: %s (密码: test123456!)\n", schoolAdmin.Username)
	fmt.Printf("  测试用户: %s (密码: test123456!)\n", testUser.Username)
}
