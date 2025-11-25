package dto

import "time"

// CreateTeamRequest 创建团队请求
type CreateTeamRequest struct {
	TeamName     string `json:"team_name" binding:"required,min=2,max=30,ascii"`
	TeamPassword string `json:"team_password" binding:"required,min=6,max=20,alphanum"`
	TeamTrack    string `json:"team_track" binding:"required,oneof=social freshman advanced"`
}

// JoinTeamRequest 加入团队请求
type JoinTeamRequest struct {
	TeamName     string `json:"team_name" binding:"required,min=2,max=30,ascii"`
	TeamPassword string `json:"team_password" binding:"required,min=6,max=20,alphanum"`
}

// UpdateTeamRequest 更新团队信息请求
type UpdateTeamRequest struct {
	TeamName     *string `json:"team_name" binding:"omitempty,min=2,max=30,ascii"`
	TeamPassword *string `json:"team_password" binding:"omitempty,min=6,max=20,alphanum"`
}

// TransferCaptainRequest 转让队长请求
type TransferCaptainRequest struct {
	NewCaptainID int64 `json:"new_captain_id" binding:"required,min=1"`
}

// RemoveMemberRequest 移除成员请求
type RemoveMemberRequest struct {
	MemberID int64 `json:"member_id" binding:"required,min=1"`
}

// UpdateTeamStatusRequest 更新团队状态请求（管理员）
type UpdateTeamStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active disbanded banned"`
}

// TeamListRequest 团队列表查询请求
type TeamListRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1,max=1000"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search    string `form:"search" binding:"omitempty,max=50"`
	TeamTrack string `form:"team_track" binding:"omitempty,oneof=social freshman advanced"`
	SchoolID  *int64 `form:"school_id"`
	Status    string `form:"status" binding:"omitempty,oneof=active disbanded banned"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=team_score created_at member_count"`
	Order     string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// TeamResponse 团队响应
type TeamResponse struct {
	ID          int64      `json:"id"`
	TeamName    string     `json:"team_name"`
	CaptainID   int64      `json:"captain_id"`
	CaptainName string     `json:"captain_name"`
	Member1ID   *int64     `json:"member1_id"`
	Member1Name *string    `json:"member1_name"`
	Member2ID   *int64     `json:"member2_id"`
	Member2Name *string    `json:"member2_name"`
	TeamTrack   string     `json:"team_track"`
	SchoolID    *int64     `json:"school_id"`
	SchoolName  *string    `json:"school_name"`
	TeamScore   int        `json:"team_score"`
	MemberCount int8       `json:"member_count"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TeamListResponse 团队列表响应
type TeamListResponse struct {
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	List  []TeamResponse `json:"list"`
}

// TeamMemberInfo 团队成员信息
type TeamMemberInfo struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"` // captain, member1, member2
}

// TeamDetailResponse 团队详情响应
type TeamDetailResponse struct {
	TeamResponse
	Members []TeamMemberInfo `json:"members"`
}

// TeamRankItem 团队排名项
type TeamRankItem struct {
	Rank        int        `json:"rank"`
	TeamID      int64      `json:"team_id"`
	TeamName    string     `json:"team_name"`
	TeamScore   int        `json:"team_score"`
	TeamTrack   string     `json:"team_track"`
	SchoolName  *string    `json:"school_name"`
	MemberCount int8       `json:"member_count"`
	SolveCount  int        `json:"solve_count"`
	LastSolveAt *time.Time `json:"last_solve_at"`
}

// TeamRankRequest 团队排名查询请求
type TeamRankRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100"`
	TeamTrack string `form:"team_track" binding:"omitempty,oneof=social freshman advanced"`
	SchoolID  *int64 `form:"school_id"`
}

// TeamRankResponse 团队排名响应
type TeamRankResponse struct {
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	List  []TeamRankItem `json:"list"`
}
