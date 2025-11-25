package routes

import (
	"isctf/controllers"
	"isctf/middleware"
	"isctf/utils"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(gin.Recovery())                // 异常恢复中间件
	r.Use(middleware.LoggerMiddleware()) // 自定义日志中间件
	r.Use(middleware.CORSMiddleware())   // 跨域中间件

	// 实例化控制器
	userController := controllers.NewUserController()
	schoolController := controllers.NewSchoolController()
	teamController := controllers.NewTeamController()

	// 健康检查接口（不需要认证）
	r.GET("/ping", func(c *gin.Context) {
		utils.Success(c, gin.H{
			"message": "pong",
			"status":  "healthy",
		})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// ---------------------------
		// 公开接口（不需要认证）
		// ---------------------------
		public := v1.Group("")
		{
			// 用户认证
			public.POST("/login", userController.Login)
			public.POST("/register/social", userController.RegisterSocial) // 社会赛道注册
			public.POST("/register/school", userController.RegisterSchool) // 院校赛道注册

			// 验证码
			public.POST("/verify/email/send", userController.SendVerifyCode)
			public.POST("/verify/email/check", userController.VerifyEmail)

			// 公开查询
			public.GET("/schools", schoolController.GetSchools)                        // 获取学校列表
			public.GET("/schools/:id", schoolController.GetSchoolByID)                 // 获取学校详情
			public.GET("/schools/name/:school_name", schoolController.GetSchoolByName) // 根据名称查学校

			// 团队公开信息
			public.GET("/teams", teamController.GetTeamList)      // 获取团队列表
			public.GET("/teams/:id", teamController.GetTeamByID)  // 获取团队详情
			public.GET("/teams/rank", teamController.GetTeamRank) // 获取团队排名
		}

		// ---------------------------
		// 需要认证的接口
		// ---------------------------
		auth := v1.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			// 用户个人中心
			users := auth.Group("/users")
			{
				users.GET("/profile", userController.GetProfile)      // 获取个人信息
				users.PUT("/profile", userController.UpdateProfile)   // 更新个人信息
				users.PUT("/password", userController.ChangePassword) // 修改密码
			}

			// 团队管理（学生端）
			teams := auth.Group("/teams")
			{
				teams.POST("", teamController.CreateTeam)               // 创建团队
				teams.POST("/join", teamController.JoinTeam)            // 加入团队
				teams.POST("/leave", teamController.LeaveTeam)          // 离开团队
				teams.GET("/me", teamController.GetMyTeam)              // 获取我的团队
				teams.PUT("", teamController.UpdateTeam)                // 更新团队信息(队长)
				teams.POST("/transfer", teamController.TransferCaptain) // 转让队长(队长)
				teams.POST("/remove", teamController.RemoveMember)      // 移除成员(队长)
				teams.POST("/disband", teamController.DisbandTeam)      // 解散团队(队长)
			}

			// ---------------------------
			// 管理员/负责人接口
			// ---------------------------

			// 学校管理员/管理员接口
			schoolAdmin := auth.Group("")
			schoolAdmin.Use(middleware.SchoolAdminMiddleware())
			{
				// 审核学生
				schoolAdmin.POST("/users/:id/verify", userController.VerifyStudent)
			}

			// 系统管理员接口
			admin := auth.Group("")
			admin.Use(middleware.AdminMiddleware())
			{
				// 学校管理
				admin.POST("/schools", schoolController.CreateSchool)
				admin.PUT("/schools/:id", schoolController.UpdateSchool)
				admin.DELETE("/schools/:id", schoolController.DeleteSchool)
				admin.PATCH("/schools/:id/status", schoolController.UpdateSchoolStatus)

				// 用户管理
				admin.GET("/users", userController.GetUserList)                   // 获取用户列表
				admin.GET("/users/:id", userController.GetUserByID)               // 获取用户详情
				admin.PATCH("/users/:id/role", userController.UpdateUserRole)     // 修改用户角色
				admin.PATCH("/users/:id/status", userController.UpdateUserStatus) // 修改用户状态

				// 团队管理
				admin.PATCH("/teams/:id/status", teamController.UpdateTeamStatus) // 管理团队状态
			}
		}
	}

	return r
}
