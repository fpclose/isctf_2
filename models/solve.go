package models

import (
	"time"
)

// Solve 解题记录模型
type Solve struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ChallengeID   int64      `gorm:"not null;index" json:"challenge_id"`
	TeamID        int64      `gorm:"not null;uniqueIndex:uk_team_challenge" json:"team_id"`
	UserID        int64      `gorm:"not null;index" json:"user_id"`
	EarnedScore   int        `gorm:"not null" json:"earned_score"`
	Rank          int        `gorm:"not null" json:"rank"`
	IsFirstBlood  bool       `gorm:"default:false;not null" json:"is_first_blood"`
	IsSecondBlood bool       `gorm:"default:false;not null" json:"is_second_blood"`
	IsThirdBlood  bool       `gorm:"default:false;not null" json:"is_third_blood"`
	SolvingTime   time.Time  `gorm:"autoCreateTime" json:"solving_time"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (Solve) TableName() string {
	return "dalictf_solve"
}
