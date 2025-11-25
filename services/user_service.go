package services

import (
	"errors"
	"fmt"
	"isctf/config"
	"isctf/dto"
	"isctf/models"
	"isctf/utils"
	"math/rand"
	"strings"
	"time"

	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct{}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{}
}

// RegisterSocial 社会赛道用户注册
func (s *UserService) RegisterSocial(req *dto.RegisterSocialRequest) (*models.User, error) {
	// 检查用户名是否已存在
	var count int64
	if err := config.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if err := config.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("邮箱已被注册")
	}

	// 创建用户
	user := &models.User{
		Username:      req.Username,
		Email:         req.Email,
		Role:          "user",
		Track:         "social",
		EmailVerified: false,
		VerifyStatus:  "pending", // 社会赛道无需审核，但保持pending状态
		Status:        "active",
	}

	// 设置密码
	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := config.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// RegisterSchool 联合院校赛道用户注册
func (s *UserService) RegisterSchool(req *dto.RegisterSchoolRequest) (*models.User, error) {
	// 检查用户名是否已存在
	var count int64
	if err := config.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if err := config.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("邮箱已被注册")
	}

	// 检查学校是否存在
	var school models.School
	if err := config.DB.First(&school, req.SchoolID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("学校不存在")
		}
		return nil, err
	}

	// 检查学校是否被封禁
	if school.Status == "suspended" {
		return nil, errors.New("该学校已被封禁，无法注册")
	}

	// 创建用户
	user := &models.User{
		Username:      req.Username,
		Email:         req.Email,
		Role:          "user",
		Track:         "school",
		SchoolID:      &req.SchoolID,
		SchoolName:    &req.SchoolName,
		UserName:      &req.UserName,
		StudentNumber: &req.StudentNumber,
		SchoolGrade:   &req.SchoolGrade,
		StudentNature: &req.StudentNature,
		EmailVerified: false,
		VerifyStatus:  "pending", // 等待院校负责人审核
		Status:        "active",
	}

	// 设置密码
	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := config.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// SendVerifyCode 发送邮箱验证码
func (s *UserService) SendVerifyCode(email string) error {
	// 生成6位数字验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 设置过期时间（5分钟）
	expiresAt := time.Now().Add(5 * time.Minute)

	// 更新用户验证码
	if err := config.DB.Model(&models.User{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"email_verify_code":      code,
			"verify_code_expires_at": expiresAt,
		}).Error; err != nil {
		return err
	}

	// TODO: 实际发送邮件
	// 这里应该调用邮件服务发送验证码
	// 暂时只打印到控制台
	fmt.Printf("📧 发送验证码到 %s: %s (有效期5分钟)\n", email, code)

	return nil
}

// VerifyEmail 验证邮箱
func (s *UserService) VerifyEmail(req *dto.VerifyEmailRequest) error {
	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 检查验证码是否正确
	if user.EmailVerifyCode == nil || *user.EmailVerifyCode != req.VerifyCode {
		return errors.New("验证码错误")
	}

	// 检查验证码是否过期
	if user.VerifyCodeExpiresAt == nil || time.Now().After(*user.VerifyCodeExpiresAt) {
		return errors.New("验证码已过期")
	}

	// 更新邮箱验证状态
	if err := config.DB.Model(&user).Updates(map[string]interface{}{
		"email_verified":         true,
		"email_verify_code":      nil,
		"verify_code_expires_at": nil,
	}).Error; err != nil {
		return err
	}

	return nil
}

// Login 用户登录
func (s *UserService) Login(req *dto.LoginRequest, ip string) (*dto.LoginResponse, error) {
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 检查密码
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status == "suspended" {
		return nil, errors.New("用户已被封禁")
	}

	// 联合院校赛道需要检查审核状态
	if user.Track == "school" && user.VerifyStatus != "approved" {
		if user.VerifyStatus == "pending" {
			return nil, errors.New("您的学生信息正在审核中，请等待院校负责人审核")
		}
		if user.VerifyStatus == "rejected" {
			return nil, errors.New("您的学生信息审核未通过，请联系院校负责人")
		}
	}

	// 更新最后登录时间和IP
	now := time.Now()
	config.DB.Model(&user).Updates(map[string]interface{}{
		"last_login_time": now,
		"last_login_ip":   ip,
	})

	// 生成JWT Token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	// 计算过期时间
	expiresAt := now.Add(24 * time.Hour).Unix()

	return &dto.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		UserInfo: dto.UserResponse{
			ID:            user.ID,
			Username:      user.Username,
			Email:         user.Email,
			Role:          user.Role,
			Track:         user.Track,
			SchoolID:      user.SchoolID,
			SchoolName:    user.SchoolName,
			UserName:      user.UserName,
			StudentNumber: user.StudentNumber,
			SchoolGrade:   user.SchoolGrade,
			StudentNature: user.StudentNature,
			EmailVerified: user.EmailVerified,
			VerifyStatus:  user.VerifyStatus,
			Status:        user.Status,
			LastLoginTime: &now,
			CreatedAt:     user.CreatedAt,
		},
	}, nil
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(id int64) (*dto.UserResponse, error) {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &dto.UserResponse{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		Role:          user.Role,
		Track:         user.Track,
		SchoolID:      user.SchoolID,
		SchoolName:    user.SchoolName,
		UserName:      user.UserName,
		StudentNumber: user.StudentNumber,
		SchoolGrade:   user.SchoolGrade,
		StudentNature: user.StudentNature,
		EmailVerified: user.EmailVerified,
		VerifyStatus:  user.VerifyStatus,
		Status:        user.Status,
		LastLoginTime: user.LastLoginTime,
		CreatedAt:     user.CreatedAt,
	}, nil
}

// UpdateProfile 更新个人信息
func (s *UserService) UpdateProfile(userID int64, req *dto.UpdateProfileRequest) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return err
	}

	updates := make(map[string]interface{})

	if req.Email != nil {
		// 检查新邮箱是否已被使用
		var count int64
		if err := config.DB.Model(&models.User{}).
			Where("email = ? AND id != ?", *req.Email, userID).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("邮箱已被使用")
		}
		updates["email"] = *req.Email
		updates["email_verified"] = false // 更换邮箱后需要重新验证
	}

	if req.UserName != nil {
		updates["user_name"] = *req.UserName
	}

	if len(updates) > 0 {
		if err := config.DB.Model(&user).Updates(updates).Error; err != nil {
			return err
		}
	}

	return nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID int64, req *dto.ChangePasswordRequest) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return err
	}

	// 验证旧密码
	if !user.CheckPassword(req.OldPassword) {
		return errors.New("旧密码错误")
	}

	// 设置新密码
	if err := user.SetPassword(req.NewPassword); err != nil {
		return err
	}

	if err := config.DB.Model(&user).Update("password", user.Password).Error; err != nil {
		return err
	}

	return nil
}

// GetUserList 获取用户列表（管理员）
func (s *UserService) GetUserList(req *dto.UserListRequest) (*dto.UserListResponse, error) {
	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Order == "" {
		req.Order = "desc"
	}

	// 构建查询
	query := config.DB.Model(&models.User{})

	// 搜索条件
	if req.Search != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR user_name LIKE ?",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	// 筛选条件
	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}
	if req.Track != "" {
		query = query.Where("track = ?", req.Track)
	}
	if req.SchoolID != nil {
		query = query.Where("school_id = ?", *req.SchoolID)
	}
	if req.VerifyStatus != "" {
		query = query.Where("verify_status = ?", req.VerifyStatus)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 排序
	sortBy := "created_at"
	if req.SortBy != "" {
		sortBy = req.SortBy
	}
	orderSQL := sortBy + " " + strings.ToUpper(req.Order)
	query = query.Order(orderSQL)

	// 分页查询
	offset := (req.Page - 1) * req.Limit
	var users []models.User
	if err := query.Offset(offset).Limit(req.Limit).Find(&users).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		list = append(list, dto.UserResponse{
			ID:            user.ID,
			Username:      user.Username,
			Email:         user.Email,
			Role:          user.Role,
			Track:         user.Track,
			SchoolID:      user.SchoolID,
			SchoolName:    user.SchoolName,
			UserName:      user.UserName,
			StudentNumber: user.StudentNumber,
			SchoolGrade:   user.SchoolGrade,
			StudentNature: user.StudentNature,
			EmailVerified: user.EmailVerified,
			VerifyStatus:  user.VerifyStatus,
			Status:        user.Status,
			LastLoginTime: user.LastLoginTime,
			CreatedAt:     user.CreatedAt,
		})
	}

	return &dto.UserListResponse{
		Total: int(total),
		Page:  req.Page,
		Limit: req.Limit,
		List:  list,
	}, nil
}

// VerifyStudent 审核学生信息（院校负责人/管理员）
func (s *UserService) VerifyStudent(userID, verifierID int64, req *dto.VerifyStudentRequest) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 检查是否是联合院校赛道
	if user.Track != "school" {
		return errors.New("该用户不是联合院校赛道，无需审核")
	}

	now := time.Now()
	updates := map[string]interface{}{
		"verify_status": req.VerifyStatus,
		"verified_by":   verifierID,
		"verified_at":   now,
	}

	if req.VerifyReason != nil {
		updates["verify_reason"] = *req.VerifyReason
	}

	// 如果被驳回，增加失败次数
	if req.VerifyStatus == "rejected" {
		updates["register_fail_count"] = user.RegisterFailCount + 1
	}

	if err := config.DB.Model(&user).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// UpdateUserRole 更新用户角色（管理员）
func (s *UserService) UpdateUserRole(userID int64, role string) error {
	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error; err != nil {
		return err
	}
	return nil
}

// UpdateUserStatus 更新用户状态（管理员）
func (s *UserService) UpdateUserStatus(userID int64, status string) error {
	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}
