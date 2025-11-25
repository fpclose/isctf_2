package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Team 团队模型
type Team struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	TeamName     string     `gorm:"type:varchar(100);not null;uniqueIndex:uk_team_name" json:"team_name"`
	TeamPassword string     `gorm:"type:varchar(255);not null" json:"-"`
	CaptainID    int64      `gorm:"not null" json:"captain_id"`
	CaptainName  string     `gorm:"type:varchar(50);not null" json:"captain_name"`
	Member1ID    *int64     `gorm:"default:null" json:"member1_id"`
	Member1Name  *string    `gorm:"type:varchar(50);default:null" json:"member1_name"`
	Member2ID    *int64     `gorm:"default:null" json:"member2_id"`
	Member2Name  *string    `gorm:"type:varchar(50);default:null" json:"member2_name"`
	TeamTrack    string     `gorm:"type:enum('social','freshman','advanced');not null" json:"team_track"`
	SchoolID     *int64     `gorm:"default:null" json:"school_id"`
	SchoolName   *string    `gorm:"type:varchar(255);default:null" json:"school_name"`
	TeamScore    int        `gorm:"default:0;not null" json:"team_score"`
	MemberCount  int8       `gorm:"default:1;not null" json:"member_count"`
	Status       string     `gorm:"type:enum('active','disbanded','banned');default:'active';not null" json:"status"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (Team) TableName() string {
	return "dalictf_team"
}

// SetPassword 设置团队密码（bcrypt加密）
func (t *Team) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	t.TeamPassword = string(hashedPassword)
	return nil
}

// CheckPassword 验证团队密码
func (t *Team) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(t.TeamPassword), []byte(password))
	return err == nil
}

// IsActive 检查团队是否处于正常状态
func (t *Team) IsActive() bool {
	return t.Status == "active"
}

// IsDisbanded 检查团队是否已解散
func (t *Team) IsDisbanded() bool {
	return t.Status == "disbanded"
}

// IsBanned 检查团队是否被封禁
func (t *Team) IsBanned() bool {
	return t.Status == "banned"
}

// IsFull 检查团队是否已满员
func (t *Team) IsFull() bool {
	return t.MemberCount >= 3
}

// IsCaptain 检查用户是否是队长
func (t *Team) IsCaptain(userID int64) bool {
	return t.CaptainID == userID
}

// IsMember 检查用户是否是团队成员
func (t *Team) IsMember(userID int64) bool {
	if t.CaptainID == userID {
		return true
	}
	if t.Member1ID != nil && *t.Member1ID == userID {
		return true
	}
	if t.Member2ID != nil && *t.Member2ID == userID {
		return true
	}
	return false
}
