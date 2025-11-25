package dto

import "time"

// RegisterSocialRequest 社会赛道注册请求
type RegisterSocialRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=50,containsany=!@#$%^&*"`
	Email    string `json:"email" binding:"required,email,max=100"`
}

// RegisterSchoolRequest 联合院校赛道注册请求
type RegisterSchoolRequest struct {
	Username      string `json:"username" binding:"required,min=3,max=20,alphanum"`
	Password      string `json:"password" binding:"required,min=8,max=50,containsany=!@#$%^&*"`
	Email         string `json:"email" binding:"required,email,max=100"`
	SchoolID      int64  `json:"school_id" binding:"required,min=1"`
	SchoolName    string `json:"school_name" binding:"required,min=2,max=100"`
	UserName      string `json:"user_name" binding:"required,min=2,max=20"`
	StudentNumber string `json:"student_number" binding:"required,min=5,max=20,alphanum"`
	SchoolGrade   string `json:"school_grade" binding:"required,min=2,max=20"`
	StudentNature string `json:"student_nature" binding:"required,oneof=undergraduate graduate"`
}

// SendVerifyCodeRequest 发送验证码请求
type SendVerifyCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyEmailRequest 验证邮箱请求
type VerifyEmailRequest struct {
	Email      string `json:"email" binding:"required,email"`
	VerifyCode string `json:"verify_code" binding:"required,len=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=1,max=50"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt int64        `json:"expires_at"`
	UserInfo  UserResponse `json:"user_info"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID                int64      `json:"id"`
	Username          string     `json:"username"`
	Email             string     `json:"email"`
	Role              string     `json:"role"`
	Track             string     `json:"track"`
	SchoolID          *int64     `json:"school_id"`
	SchoolName        *string    `json:"school_name"`
	UserName          *string    `json:"user_name"`
	StudentNumber     *string    `json:"student_number"`
	SchoolGrade       *string    `json:"school_grade"`
	StudentNature     *string    `json:"student_nature"`
	EmailVerified     bool       `json:"email_verified"`
	VerifyStatus      string     `json:"verify_status"`
	Status            string     `json:"status"`
	LastLoginTime     *time.Time `json:"last_login_time"`
	CreatedAt         time.Time  `json:"created_at"`
}

// UpdateProfileRequest 更新个人信息请求
type UpdateProfileRequest struct {
	Email    *string `json:"email" binding:"omitempty,email,max=100"`
	UserName *string `json:"user_name" binding:"omitempty,min=2,max=20"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=1,max=50"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=50,containsany=!@#$%^&*"`
}

// VerifyStudentRequest 审核学生信息请求（院校负责人）
type VerifyStudentRequest struct {
	VerifyStatus string  `json:"verify_status" binding:"required,oneof=approved rejected"`
	VerifyReason *string `json:"verify_reason"`
}

// UserListRequest 用户列表查询请求
type UserListRequest struct {
	Page         int    `form:"page" binding:"omitempty,min=1"`
	Limit        int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search       string `form:"search"`
	Role         string `form:"role" binding:"omitempty,oneof=user school_admin admin super_admin"`
	Track        string `form:"track" binding:"omitempty,oneof=social school"`
	SchoolID     *int64 `form:"school_id"`
	VerifyStatus string `form:"verify_status" binding:"omitempty,oneof=pending approved rejected"`
	Status       string `form:"status" binding:"omitempty,oneof=active suspended"`
	SortBy       string `form:"sort_by" binding:"omitempty,oneof=created_at last_login_time"`
	Order        string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	List  []UserResponse `json:"list"`
}

// UpdateUserRoleRequest 更新用户角色请求（管理员）
type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user school_admin admin super_admin"`
}

// UpdateUserStatusRequest 更新用户状态请求（管理员）
type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active suspended"`
}
