package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
type User struct {
	ID                  int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username            string     `gorm:"type:varchar(50);not null;uniqueIndex:uk_username" json:"username"`
	Password            string     `gorm:"type:varchar(255);not null" json:"-"`
	Email               string     `gorm:"type:varchar(100);not null;uniqueIndex:uk_email" json:"email"`
	Role                string     `gorm:"type:enum('user','school_admin','admin','super_admin');default:'user';not null" json:"role"`
	Track               string     `gorm:"type:enum('social','school');default:'social';not null" json:"track"`
	SchoolID            *int64     `gorm:"default:null" json:"school_id"`
	SchoolName          *string    `gorm:"type:varchar(255);default:null" json:"school_name"`
	UserName            *string    `gorm:"type:varchar(50);default:null" json:"user_name"`
	StudentNumber       *string    `gorm:"type:varchar(50);default:null" json:"student_number"`
	SchoolGrade         *string    `gorm:"type:varchar(10);default:null" json:"school_grade"`
	StudentNature       *string    `gorm:"type:enum('undergraduate','graduate');default:null" json:"student_nature"`
	EmailVerified       bool       `gorm:"default:0;not null" json:"email_verified"`
	EmailVerifyCode     *string    `gorm:"type:varchar(10);default:null" json:"-"`
	VerifyCodeExpiresAt *time.Time `gorm:"default:null" json:"-"`
	RegisterFailCount   int        `gorm:"default:0;not null" json:"register_fail_count"`
	VerifyStatus        string     `gorm:"type:enum('pending','approved','rejected');default:'pending';not null" json:"verify_status"`
	VerifyReason        *string    `gorm:"type:varchar(500);default:null" json:"verify_reason"`
	VerifiedBy          *int64     `gorm:"default:null" json:"verified_by"`
	VerifiedAt          *time.Time `gorm:"default:null" json:"verified_at"`
	Status              string     `gorm:"type:enum('active','suspended');default:'active';not null" json:"status"`
	LastLoginTime       *time.Time `gorm:"default:null" json:"last_login_time"`
	LastLoginIP         *string    `gorm:"type:varchar(50);default:null" json:"last_login_ip"`
	CreatedAt           time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt           *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "dalictf_user"
}

// SetPassword 设置密码（bcrypt加密）
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsActive 检查用户是否处于正常状态
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// IsSuspended 检查用户是否被封禁
func (u *User) IsSuspended() bool {
	return u.Status == "suspended"
}

// IsAdmin 检查是否是管理员
func (u *User) IsAdmin() bool {
	return u.Role == "admin" || u.Role == "super_admin"
}

// IsSchoolAdmin 检查是否是学校管理员
func (u *User) IsSchoolAdmin() bool {
	return u.Role == "school_admin"
}

// IsSocialTrack 检查是否是社会赛道
func (u *User) IsSocialTrack() bool {
	return u.Track == "social"
}

// IsSchoolTrack 检查是否是联合院校赛道
func (u *User) IsSchoolTrack() bool {
	return u.Track == "school"
}

// IsApproved 检查学生信息是否已审核通过
func (u *User) IsApproved() bool {
	return u.VerifyStatus == "approved"
}
