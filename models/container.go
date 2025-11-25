package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// PortMapping 端口映射 JSON 结构
type PortMapping map[string]string

func (p PortMapping) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *PortMapping) Scan(value interface{}) error {
	if value == nil {
		*p = make(PortMapping)
		return nil
	}
	bytes, _ := value.([]byte)
	return json.Unmarshal(bytes, p)
}

// Container 容器实例模型
type Container struct {
	ID            int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	ChallengeID   int64       `gorm:"not null;index" json:"challenge_id"`
	TeamID        int64       `gorm:"not null;index" json:"team_id"`
	UserID        int64       `gorm:"not null;index" json:"user_id"`
	ContainerName string      `gorm:"type:varchar(255);not null" json:"container_name"` // 存储 Docker ID
	DockerImage   string      `gorm:"type:varchar(255);not null" json:"docker_image"`
	DockerPorts   DockerPorts `gorm:"type:json" json:"docker_ports"` // 复用 challenge.go 中的定义
	HostMapping   PortMapping `gorm:"type:json" json:"host_mapping"`
	ContainerFlag string      `gorm:"type:varchar(500);not null" json:"-"` // Flag 不直接返回给前端
	State         string      `gorm:"type:enum('running','stopped','destroyed');default:'running';not null" json:"state"`
	StartTime     time.Time   `gorm:"autoCreateTime" json:"start_time"`
	EndTime       time.Time   `gorm:"not null" json:"end_time"`
	ExtendedCount int8        `gorm:"default:0;not null" json:"extended_count"`
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     *time.Time  `gorm:"index" json:"deleted_at,omitempty"`
}

func (Container) TableName() string {
	return "dalictf_container"
}
