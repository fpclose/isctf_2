package controllers

import (
	"isctf/dto"
	"isctf/services"
	"isctf/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService *services.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// RegisterSocial 社会赛道用户注册
func (c *UserController) RegisterSocial(ctx *gin.Context) {
	var req dto.RegisterSocialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	user, err := c.userService.RegisterSocial(&req)
	if err != nil {
		if err.Error() == "用户名已存在" || err.Error() == "邮箱已被注册" {
			utils.ErrorWithMsg(ctx, utils.USER_ALREADY_EXIST, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "注册失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "注册成功", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"track":    user.Track,
	})
}

// RegisterSchool 联合院校赛道用户注册
func (c *UserController) RegisterSchool(ctx *gin.Context) {
	var req dto.RegisterSchoolRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	user, err := c.userService.RegisterSchool(&req)
	if err != nil {
		if err.Error() == "用户名已存在" || err.Error() == "邮箱已被注册" {
			utils.ErrorWithMsg(ctx, utils.USER_ALREADY_EXIST, err.Error())
			return
		}
		if err.Error() == "学校不存在" {
			utils.ErrorWithMsg(ctx, utils.SCHOOL_NOT_EXIST, err.Error())
			return
		}
		if err.Error() == "该学校已被封禁，无法注册" {
			utils.ErrorWithMsg(ctx, utils.SCHOOL_SUSPENDED, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "注册失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "注册成功，请等待院校负责人审核", gin.H{
		"id":            user.ID,
		"username":      user.Username,
		"email":         user.Email,
		"role":          user.Role,
		"track":         user.Track,
		"verify_status": user.VerifyStatus,
	})
}

// SendVerifyCode 发送邮箱验证码
func (c *UserController) SendVerifyCode(ctx *gin.Context) {
	var req dto.SendVerifyCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	if err := c.userService.SendVerifyCode(req.Email); err != nil {
		utils.ErrorWithMsg(ctx, utils.ERROR, "发送验证码失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "验证码已发送，请查收邮件", nil)
}

// VerifyEmail 验证邮箱
func (c *UserController) VerifyEmail(ctx *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	if err := c.userService.VerifyEmail(&req); err != nil {
		if err.Error() == "用户不存在" {
			utils.ErrorWithMsg(ctx, utils.USER_NOT_EXIST, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "邮箱验证成功", nil)
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	// 获取客户端IP
	ip := ctx.ClientIP()

	result, err := c.userService.Login(&req, ip)
	if err != nil {
		if err.Error() == "用户名或密码错误" {
			utils.ErrorWithMsg(ctx, utils.PASSWORD_ERROR, err.Error())
			return
		}
		if err.Error() == "用户已被封禁" {
			utils.ErrorWithMsg(ctx, utils.FORBIDDEN, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "登录成功", result)
}

// GetProfile 获取个人信息
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Error(ctx, utils.UNAUTHORIZED)
		return
	}

	userInfo, err := c.userService.GetUserByID(userID.(int64))
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.ERROR, "获取用户信息失败: "+err.Error())
		return
	}

	utils.Success(ctx, userInfo)
}

// UpdateProfile 更新个人信息
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Error(ctx, utils.UNAUTHORIZED)
		return
	}

	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	if err := c.userService.UpdateProfile(userID.(int64), &req); err != nil {
		if err.Error() == "邮箱已被使用" {
			utils.ErrorWithMsg(ctx, utils.USER_ALREADY_EXIST, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "更新信息失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "更新信息成功", nil)
}

// ChangePassword 修改密码
func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Error(ctx, utils.UNAUTHORIZED)
		return
	}

	var req dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	if err := c.userService.ChangePassword(userID.(int64), &req); err != nil {
		if err.Error() == "旧密码错误" {
			utils.ErrorWithMsg(ctx, utils.PASSWORD_ERROR, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "修改密码失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "修改密码成功", nil)
}

// GetUserByID 根据ID获取用户信息（管理员）
func (c *UserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "无效的用户ID")
		return
	}

	userInfo, err := c.userService.GetUserByID(id)
	if err != nil {
		if err.Error() == "用户不存在" {
			utils.Error(ctx, utils.USER_NOT_EXIST)
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "获取用户信息失败")
		return
	}

	utils.Success(ctx, userInfo)
}

// GetUserList 获取用户列表（管理员）
func (c *UserController) GetUserList(ctx *gin.Context) {
	var req dto.UserListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	result, err := c.userService.GetUserList(&req)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.ERROR, "获取用户列表失败: "+err.Error())
		return
	}

	utils.Success(ctx, result)
}

// VerifyStudent 审核学生信息（院校负责人/管理员）
func (c *UserController) VerifyStudent(ctx *gin.Context) {
	verifierID, exists := ctx.Get("user_id")
	if !exists {
		utils.Error(ctx, utils.UNAUTHORIZED)
		return
	}

	idStr := ctx.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "无效的用户ID")
		return
	}

	var req dto.VerifyStudentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	if err := c.userService.VerifyStudent(userID, verifierID.(int64), &req); err != nil {
		if err.Error() == "用户不存在" {
			utils.ErrorWithMsg(ctx, utils.USER_NOT_EXIST, err.Error())
			return
		}
		if err.Error() == "该用户不是联合院校赛道，无需审核" {
			utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "审核失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "审核成功", nil)
}

// UpdateUserRole 更新用户角色（管理员）
func (c *UserController) UpdateUserRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "无效的用户ID")
		return
	}

	var req dto.UpdateUserRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	if err := c.userService.UpdateUserRole(userID, req.Role); err != nil {
		utils.ErrorWithMsg(ctx, utils.ERROR, "更新角色失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "更新角色成功", nil)
}

// UpdateUserStatus 更新用户状态（管理员）
func (c *UserController) UpdateUserStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "无效的用户ID")
		return
	}

	var req dto.UpdateUserStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	if err := c.userService.UpdateUserStatus(userID, req.Status); err != nil {
		utils.ErrorWithMsg(ctx, utils.ERROR, "更新状态失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "更新状态成功", nil)
}
