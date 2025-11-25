package models

import (
	"time"
)

// School 参赛学校模型
type School struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	SchoolName  string     `gorm:"type:varchar(255);not null;uniqueIndex:uk_school_name" json:"school_name"`
	SchoolAdmin *int64     `gorm:"default:null" json:"school_admin"`
	UserCount   int        `gorm:"default:0;not null" json:"user_count"`
	Status      string     `gorm:"type:enum('active','suspended');default:'active';not null" json:"status"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (School) TableName() string {
	return "dalictf_school"
}

// IsActive 检查学校是否处于正常状态
func (s *School) IsActive() bool {
	return s.Status == "active"
}

// IsSuspended 检查学校是否被封禁
func (s *School) IsSuspended() bool {
	return s.Status == "suspended"
}
