package dto

import "time"

// CategoryRequest 题目类型请求
type CategoryRequest struct {
	Direction   string  `json:"direction" binding:"required,max=50"`
	NameZh      string  `json:"name_zh" binding:"required,max=50"`
	NameEn      string  `json:"name_en" binding:"required,max=50"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	Color       *string `json:"color"`
	SortOrder   int     `json:"sort_order"`
}

// ChallengeRequest 题目创建/更新请求
type ChallengeRequest struct {
	ChallengeName string            `json:"challenge_name" binding:"required"`
	Direction     string            `json:"direction" binding:"required"`
	Author        string            `json:"author" binding:"required"`
	Description   string            `json:"description" binding:"required"`
	Hint          *string           `json:"hint"`
	State         string            `json:"state" binding:"oneof=visible hidden"`
	Mode          string            `json:"mode" binding:"oneof=static dynamic"`
	StaticFlag    *string           `json:"static_flag"`
	DockerImage   *string           `json:"docker_image"`
	DockerPorts   map[string]string `json:"docker_ports"`
	Difficulty    string            `json:"difficulty" binding:"oneof=easy medium hard expert"`
	InitialScore  int               `json:"initial_score" binding:"min=1"`
	MinScore      int               `json:"min_score" binding:"min=0"`
	DecayRatio    float64           `json:"decay_ratio" binding:"min=0.1,max=1.0"`
}

// SubmitFlagRequest 提交 Flag 请求
type SubmitFlagRequest struct {
	Flag string `json:"flag" binding:"required"`
}

// ContainerInfo 容器信息响应
type ContainerInfo struct {
	ContainerID   string            `json:"container_id"`
	ChallengeID   int64             `json:"challenge_id"`
	Status        string            `json:"status"`
	Host          string            `json:"host"`
	Ports         map[string]string `json:"ports"`
	ExpiresAt     time.Time         `json:"expires_at"`
	TimeRemaining string            `json:"time_remaining"`
}

// ChallengeListRequest 题目列表查询参数
type ChallengeListRequest struct {
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
	Direction  string `form:"direction"`
	Difficulty string `form:"difficulty"`
	Search     string `form:"search"`
	State      string `form:"state"` // 管理员用
}

// LogListRequest 日志查询参数
type LogListRequest struct {
	Page        int    `form:"page"`
	Limit       int    `form:"limit"`
	ChallengeID int64  `form:"challenge_id"`
	TeamID      int64  `form:"team_id"`
	UserID      int64  `form:"user_id"`
	FlagResult  string `form:"flag_result"`
	Search      string `form:"search"`
}
