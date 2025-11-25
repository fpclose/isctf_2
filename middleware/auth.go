package middleware

import (
	"isctf/config"
	"isctf/models"
	"isctf/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization 头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, utils.UNAUTHORIZED)
			c.Abort()
			return
		}

		// 检查 token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.ErrorWithMsg(c, utils.TOKEN_INVALID, "Token格式错误")
			c.Abort()
			return
		}

		// 解析 token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.ErrorWithMsg(c, utils.TOKEN_INVALID, "Token无效或已过期")
			c.Abort()
			return
		}

		// 验证用户是否存在且处于正常状态
		var user models.User
		if err := config.DB.Where("id = ? AND status = ?", claims.UserID, "active").First(&user).Error; err != nil {
			utils.ErrorWithMsg(c, utils.TOKEN_INVALID, "用户不存在或已被封禁")
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 确保用户已通过认证
		userID, exists := c.Get("user_id")
		if !exists {
			utils.Error(c, utils.UNAUTHORIZED)
			c.Abort()
			return
		}

		// 获取用户角色
		role, exists := c.Get("role")
		if !exists {
			utils.Error(c, utils.PERMISSION_DENIED)
			c.Abort()
			return
		}

		// 验证管理员权限
		if role != "admin" && role != "super_admin" {
			utils.ErrorWithMsg(c, utils.PERMISSION_DENIED, "需要管理员权限")
			c.Abort()
			return
		}

		// 再次验证数据库中的用户角色
		var user models.User
		if err := config.DB.Where("id = ? AND role IN ?", userID, []string{"admin", "super_admin"}).First(&user).Error; err != nil {
			utils.ErrorWithMsg(c, utils.PERMISSION_DENIED, "用户权限验证失败")
			c.Abort()
			return
		}

		c.Next()
	}
}

// SchoolAdminMiddleware 学校管理员权限中间件
func SchoolAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 确保用户已通过认证
		userID, exists := c.Get("user_id")
		if !exists {
			utils.Error(c, utils.UNAUTHORIZED)
			c.Abort()
			return
		}

		// 获取用户角色
		role, exists := c.Get("role")
		if !exists {
			utils.Error(c, utils.PERMISSION_DENIED)
			c.Abort()
			return
		}

		// 验证学校管理员权限（admin, school_admin, super_admin）
		if role != "admin" && role != "school_admin" && role != "super_admin" {
			utils.ErrorWithMsg(c, utils.PERMISSION_DENIED, "需要学校管理员权限")
			c.Abort()
			return
		}

		// 再次验证数据库中的用户角色
		var user models.User
		if err := config.DB.Where("id = ? AND role IN ?", userID, []string{"admin", "school_admin", "super_admin"}).First(&user).Error; err != nil {
			utils.ErrorWithMsg(c, utils.PERMISSION_DENIED, "用户权限验证失败")
			c.Abort()
			return
		}

		c.Next()
	}
}
