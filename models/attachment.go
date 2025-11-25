package models

import "time"

// Attachment 题目附件模型
type Attachment struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ChallengeID  int64      `gorm:"not null;index" json:"challenge_id"`
	Storage      string     `gorm:"type:enum('url','object');default:'object';not null" json:"storage"`
	URL          *string    `gorm:"type:varchar(1000);default:null" json:"url"`
	ObjectBucket *string    `gorm:"type:varchar(100);default:null" json:"object_bucket"`
	ObjectKey    *string    `gorm:"type:varchar(500);default:null" json:"object_key"`
	FileName     string     `gorm:"type:varchar(255);not null" json:"file_name"`
	ContentType  *string    `gorm:"type:varchar(100);default:null" json:"content_type"`
	FileSize     *int64     `gorm:"default:null" json:"file_size"`
	SHA256       *string    `gorm:"type:char(64);default:null" json:"sha256"`
	Status       string     `gorm:"type:enum('pending','active','infected','error');default:'pending';not null" json:"status"`
	Visibility   string     `gorm:"type:enum('public','private','team');default:'private';not null" json:"visibility"`
	Version      string     `gorm:"type:varchar(50);default:'1.0'" json:"version"`
	SortOrder    int        `gorm:"default:0;not null" json:"sort_order"`
	CreatedBy    *int64     `gorm:"default:null" json:"created_by"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Attachment) TableName() string {
	return "dalictf_challenge_attachment"
}
