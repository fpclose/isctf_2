package models

import "time"

// SubmissionLog 提交日志模型
type SubmissionLog struct {
	ID             int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ChallengeID    int64      `gorm:"not null;index" json:"challenge_id"`
	TeamID         int64      `gorm:"not null;index" json:"team_id"`
	UserID         int64      `gorm:"not null;index" json:"user_id"`
	SubmittedFlag  string     `gorm:"type:varchar(500);not null" json:"submitted_flag"`
	FlagResult     string     `gorm:"type:enum('correct','wrong','duplicate');not null" json:"flag_result"`
	ChallengeType  string     `gorm:"type:enum('static','dynamic');not null" json:"challenge_type"`
	IPAddress      string     `gorm:"type:varchar(50);not null" json:"ip_address"`
	UserAgent      *string    `gorm:"type:varchar(500)" json:"user_agent"`
	SubmissionTime time.Time  `gorm:"autoCreateTime" json:"submission_time"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (SubmissionLog) TableName() string {
	return "dalictf_submission_log"
}
