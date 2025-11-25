package dto

import "time"

// CreateSchoolRequest 创建学校请求
type CreateSchoolRequest struct {
	SchoolName  string `json:"school_name" binding:"required,min=2,max=100,ascii"`
	SchoolAdmin *int64 `json:"school_admin" binding:"omitempty,min=1"`
}

// UpdateSchoolRequest 更新学校请求
type UpdateSchoolRequest struct {
	SchoolName  *string `json:"school_name" binding:"omitempty,min=2,max=100,ascii"`
	SchoolAdmin *int64  `json:"school_admin" binding:"omitempty,min=1"`
}

// UpdateSchoolStatusRequest 更新学校状态请求
type UpdateSchoolStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active suspended"`
}

// SchoolListRequest 学校列表查询请求
type SchoolListRequest struct {
	Page    int    `form:"page" binding:"omitempty,min=1,max=1000"`
	Limit   int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search  string `form:"search" binding:"omitempty,max=50"`
	SortBy  string `form:"sort_by" binding:"omitempty,oneof=user_count created_at"`
	Order   string `form:"order" binding:"omitempty,oneof=asc desc"`
	Status  string `form:"status" binding:"omitempty,oneof=active suspended"`
}

// SchoolResponse 学校响应
type SchoolResponse struct {
	ID          int64     `json:"id"`
	SchoolName  string    `json:"school_name"`
	SchoolAdmin *int64    `json:"school_admin"`
	AdminName   *string   `json:"admin_name,omitempty"`
	UserCount   int       `json:"user_count"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SchoolListResponse 学校列表响应
type SchoolListResponse struct {
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
	List  []SchoolResponse `json:"list"`
}
