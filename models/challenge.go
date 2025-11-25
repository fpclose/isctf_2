package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// DockerPorts Docker端口映射配置
type DockerPorts map[string]string

// Value 实现driver.Valuer接口
func (dp DockerPorts) Value() (driver.Value, error) {
	return json.Marshal(dp)
}

// Scan 实现sql.Scanner接口
func (dp *DockerPorts) Scan(value interface{}) error {
	if value == nil {
		*dp = make(DockerPorts)
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, dp)
	case string:
		return json.Unmarshal([]byte(v), dp)
	}
	return nil
}

// Challenge 题目模型
type Challenge struct {
	ID            int64        `json:"id" gorm:"primaryKey;autoIncrement;comment:题目主键ID"`
	ChallengeName string       `json:"challenge_name" gorm:"type:varchar(255);not null;comment:题目名称"`
	Direction     string       `json:"direction" gorm:"type:varchar(50);not null;index:idx_direction;comment:题目类型"`
	Author        string       `json:"author" gorm:"type:varchar(100);not null;comment:出题人网名"`
	Description   string       `json:"description" gorm:"type:text;not null;comment:题目描述"`
	Hint          *string      `json:"hint" gorm:"type:text;comment:题目提示"`
	State         string       `json:"state" gorm:"type:enum('visible','hidden');not null;default:'visible';index:idx_state;comment:题目状态"`
	Mode          string       `json:"mode" gorm:"type:enum('static','dynamic');not null;default:'static';index:idx_mode;comment:题目模式"`
	StaticFlag    *string      `json:"static_flag" gorm:"type:varchar(500);comment:静态题flag"`
	DockerImage   *string      `json:"docker_image" gorm:"type:varchar(255);comment:动态题Docker镜像"`
	DockerPorts   DockerPorts  `json:"docker_ports" gorm:"type:json;comment:容器端口映射"`
	Difficulty    string       `json:"difficulty" gorm:"type:enum('easy','medium','hard','expert');not null;default:'medium';index:idx_difficulty;comment:题目难度"`
	InitialScore  int          `json:"initial_score" gorm:"not null;default:100;comment:初始分值"`
	MinScore      int          `json:"min_score" gorm:"not null;default:50;comment:最低分值"`
	CurrentScore  int          `json:"current_score" gorm:"not null;default:100;index:idx_current_score;comment:当前分值"`
	DecayRatio    float64      `json:"decay_ratio" gorm:"type:decimal(5,2);not null;default:0.90;comment:分数衰减比率"`
	SolvedCount   int          `json:"solved_count" gorm:"not null;default:0;index:idx_solved_count;comment:解出次数"`
	CreatedAt     time.Time    `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_created_at;comment:创建时间"`
	UpdatedAt     time.Time    `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
	DeletedAt     *time.Time   `json:"deleted_at" gorm:"index;comment:软删除时间"`
	
	// 关联字段
	Category      *ChallengeCategory `json:"category,omitempty" gorm:"foreignKey:Direction;references:Direction"`
}

// TableName 指定表名
func (Challenge) TableName() string {
	return "dalictf_challenge"
}

// IsVisible 检查题目是否可见
func (c *Challenge) IsVisible() bool {
	return c.State == "visible" && c.DeletedAt == nil
}

// IsStatic 检查是否为静态题目
func (c *Challenge) IsStatic() bool {
	return c.Mode == "static"
}

// IsDynamic 检查是否为动态题目
func (c *Challenge) IsDynamic() bool {
	return c.Mode == "dynamic"
}

// GetDifficultyLevel 获取难度等级（数字）
func (c *Challenge) GetDifficultyLevel() int {
	switch c.Difficulty {
	case "easy":
		return 1
	case "medium":
		return 2
	case "hard":
		return 3
	case "expert":
		return 4
	default:
		return 2
	}
}

// GetDifficultyColor 获取难度对应的颜色
func (c *Challenge) GetDifficultyColor() string {
	switch c.Difficulty {
	case "easy":
		return "#4CAF50"    // 绿色
	case "medium":
		return "#FF9800"    // 橙色
	case "hard":
		return "#F44336"    // 红色
	case "expert":
		return "#9C27B0"    // 紫色
	default:
		return "#9E9E9E"    // 灰色
	}
}

// CalculateCurrentScore 计算当前分数（基于解出次数）
func (c *Challenge) CalculateCurrentScore() int {
	if c.SolvedCount == 0 {
		return c.InitialScore
	}
	
	// 使用衰减公式：current_score = max(initial_score * decay_ratio^solved_count, min_score)
	score := float64(c.InitialScore)
	for i := 0; i < c.SolvedCount; i++ {
		score *= c.DecayRatio
	}
	
	currentScore := int(score)
	if currentScore < c.MinScore {
		return c.MinScore
	}
	return currentScore
}

// UpdateScore 更新题目分数（在有人解出后调用）
func (c *Challenge) UpdateScore(tx *gorm.DB) error {
	c.SolvedCount++
	c.CurrentScore = c.CalculateCurrentScore()
	return tx.Save(c).Error
}

// BeforeCreate GORM钩子 - 创建前
func (c *Challenge) BeforeCreate(tx *gorm.DB) error {
	// 验证静态题目必须有flag
	if c.Mode == "static" && (c.StaticFlag == nil || *c.StaticFlag == "") {
		return gorm.ErrInvalidField
	}
	
	// 验证动态题目必须有Docker镜像
	if c.Mode == "dynamic" && (c.DockerImage == nil || *c.DockerImage == "") {
		return gorm.ErrInvalidField
	}
	
	// 初始化当前分数
	if c.CurrentScore == 0 {
		c.CurrentScore = c.InitialScore
	}
	
	return nil
}

// BeforeUpdate GORM钩子 - 更新前
func (c *Challenge) BeforeUpdate(tx *gorm.DB) error {
	// 如果题目模式改变，验证相应字段
	if tx.Statement.Changed("Mode") {
		if c.Mode == "static" && (c.StaticFlag == nil || *c.StaticFlag == "") {
			return gorm.ErrInvalidField
		}
		if c.Mode == "dynamic" && (c.DockerImage == nil || *c.DockerImage == "") {
			return gorm.ErrInvalidField
		}
	}
	
	return nil
}

// ChallengeStatistics 题目统计信息
type ChallengeStatistics struct {
	TotalChallenges    int `json:"total_challenges"`
	VisibleChallenges  int `json:"visible_challenges"`
	HiddenChallenges   int `json:"hidden_challenges"`
	StaticChallenges   int `json:"static_challenges"`
	DynamicChallenges  int `json:"dynamic_challenges"`
	TotalSolves        int `json:"total_solves"`
	AverageScore       float64 `json:"average_score"`
	HardestChallenge   *Challenge `json:"hardest_challenge,omitempty"`
	EasiestChallenge   *Challenge `json:"easiest_challenge,omitempty"`
}
