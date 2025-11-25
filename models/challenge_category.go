package models

import (
	"time"
)

// ChallengeCategory 题目类型分类模型
type ChallengeCategory struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:分类主键ID"`
	Direction   string    `json:"direction" gorm:"type:varchar(50);not null;uniqueIndex:uk_direction;comment:类型方向"`
	NameZh      string    `json:"name_zh" gorm:"type:varchar(50);not null;comment:类型中文名称"`
	NameEn      string    `json:"name_en" gorm:"type:varchar(50);not null;comment:类型英文名称"`
	Description *string   `json:"description" gorm:"type:varchar(500);comment:类型描述"`
	Icon        *string   `json:"icon" gorm:"type:varchar(100);comment:图标标识"`
	Color       *string   `json:"color" gorm:"type:varchar(20);comment:主题颜色"`
	SortOrder   int       `json:"sort_order" gorm:"not null;default:0;comment:排序顺序"`
	Status      string    `json:"status" gorm:"type:enum('active','inactive');not null;default:'active';comment:分类状态"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index;comment:软删除时间"`
}

// TableName 指定表名
func (ChallengeCategory) TableName() string {
	return "dalictf_challenge_category"
}

// IsActive 检查分类是否启用
func (c *ChallengeCategory) IsActive() bool {
	return c.Status == "active"
}

// IsActiveAndNotDeleted 检查分类是否启用且未删除
func (c *ChallengeCategory) IsActiveAndNotDeleted() bool {
	return c.IsActive() && c.DeletedAt == nil
}

// GetDisplayName 获取显示名称（根据语言偏好）
func (c *ChallengeCategory) GetDisplayName(lang string) string {
	if lang == "en" && c.NameEn != "" {
		return c.NameEn
	}
	return c.NameZh
}

// GetDefaultCategories 获取默认分类数据
func GetDefaultCategories() []ChallengeCategory {
	return []ChallengeCategory{
		{
			Direction:   "Web",
			NameZh:      "Web安全",
			NameEn:      "Web Security",
			Description: stringPtr("Web应用安全，包括SQL注入、XSS、CSRF等常见Web漏洞"),
			Icon:        stringPtr("icon-web"),
			Color:       stringPtr("#FF6B6B"),
			SortOrder:   1,
			Status:      "active",
		},
		{
			Direction:   "Misc",
			NameZh:      "杂项",
			NameEn:      "Miscellaneous",
			Description: stringPtr("杂项题目，包括编码、隐写、社工等多种技巧"),
			Icon:        stringPtr("icon-misc"),
			Color:       stringPtr("#4ECDC4"),
			SortOrder:   2,
			Status:      "active",
		},
		{
			Direction:   "Crypto",
			NameZh:      "密码学",
			NameEn:      "Cryptography",
			Description: stringPtr("密码学相关题目，包括古典密码、现代加密算法等"),
			Icon:        stringPtr("icon-crypto"),
			Color:       stringPtr("#FFE66D"),
			SortOrder:   3,
			Status:      "active",
		},
		{
			Direction:   "Reverse",
			NameZh:      "逆向工程",
			NameEn:      "Reverse Engineering",
			Description: stringPtr("二进制程序逆向分析，包括软件破解、协议分析等"),
			Icon:        stringPtr("icon-reverse"),
			Color:       stringPtr("#A8E6CF"),
			SortOrder:   4,
			Status:      "active",
		},
		{
			Direction:   "Pwn",
			NameZh:      "二进制漏洞",
			NameEn:      "Binary Exploitation",
			Description: stringPtr("二进制漏洞利用，包括栈溢出、堆漏洞等"),
			Icon:        stringPtr("icon-pwn"),
			Color:       stringPtr("#FF8B94"),
			SortOrder:   5,
			Status:      "active",
		},
		{
			Direction:   "Forensics",
			NameZh:      "取证分析",
			NameEn:      "Forensics",
			Description: stringPtr("数字取证分析，包括文件恢复、内存分析、网络流量分析等"),
			Icon:        stringPtr("icon-forensics"),
			Color:       stringPtr("#B4A7D6"),
			SortOrder:   6,
			Status:      "active",
		},
		{
			Direction:   "Mobile",
			NameZh:      "移动安全",
			NameEn:      "Mobile Security",
			Description: stringPtr("移动应用安全，包括Android、iOS应用逆向和漏洞分析"),
			Icon:        stringPtr("icon-mobile"),
			Color:       stringPtr("#FFB347"),
			SortOrder:   7,
			Status:      "active",
		},
		{
			Direction:   "Blockchain",
			NameZh:      "区块链安全",
			NameEn:      "Blockchain Security",
			Description: stringPtr("区块链安全，包括智能合约审计、DeFi漏洞分析等"),
			Icon:        stringPtr("icon-blockchain"),
			Color:       stringPtr("#87CEEB"),
			SortOrder:   8,
			Status:      "active",
		},
		{
			Direction:   "IoT",
			NameZh:      "物联网安全",
			NameEn:      "IoT Security",
			Description: stringPtr("物联网设备安全，包括固件分析、硬件漏洞等"),
			Icon:        stringPtr("icon-iot"),
			Color:       stringPtr("#98D8C8"),
			SortOrder:   9,
			Status:      "active",
		},
		{
			Direction:   "AI",
			NameZh:      "AI安全",
			NameEn:      "AI Security",
			Description: stringPtr("人工智能安全，包括模型对抗攻击、数据隐私保护等"),
			Icon:        stringPtr("icon-ai"),
			Color:       stringPtr("#DDA0DD"),
			SortOrder:   10,
			Status:      "active",
		},
	}
}

// stringPtr 返回字符串指针的辅助函数
func stringPtr(s string) *string {
	return &s
}
