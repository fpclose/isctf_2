package main

import (
	"fmt"
	"isctf/cmd"
	"isctf/config"
	"isctf/routes"
	"os"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("ISCTF 平台启动中...")
	fmt.Println("========================================")

	config.InitConfig()
	config.ConnectDatabase()

	// 检查是否需要初始化管理员账户
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			fmt.Println("正在初始化系统数据...")
			cmd.InitDefaultAdmin()
			cmd.InitTestData()
			fmt.Println("初始化完成！")
			return
		case "reset-admin":
			fmt.Println("正在重置管理员密码...")
			cmd.ResetAdminPassword()
			fmt.Println("重置完成！")
			return
		case "recreate-admin":
			fmt.Println("正在重新创建管理员账户...")
			cmd.DeleteAndRecreateAdmin()
			fmt.Println("重新创建完成！")
			return
		}
	}

	// 自动检查并创建默认管理员
	cmd.InitDefaultAdmin()

	r := routes.SetupRouter()

	serverAddr := ":" + config.AppConfig.Server.Port
	fmt.Printf("\n服务器启动成功，监听地址: http://localhost%s\n", serverAddr)
	fmt.Println("========================================")
	fmt.Println("默认管理员账户:")
	fmt.Println("  用户名: root_admin")
	fmt.Println("  密码: fpclose_SfTian_i5ctf")
	fmt.Println("========================================")
	
	err := r.Run(serverAddr)
	if err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
		return
	}
}
