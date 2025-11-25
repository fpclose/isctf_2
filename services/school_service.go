package services

import (
	"errors"
	"isctf/config"
	"isctf/dto"
	"isctf/models"
	"strings"

	"gorm.io/gorm"
)

// SchoolService 学校服务
type SchoolService struct{}

// NewSchoolService 创建学校服务实例
func NewSchoolService() *SchoolService {
	return &SchoolService{}
}

// CreateSchool 创建学校
func (s *SchoolService) CreateSchool(req *dto.CreateSchoolRequest) (*models.School, error) {
	// 检查学校名称是否已存在
	var count int64
	if err := config.DB.Model(&models.School{}).Where("school_name = ?", req.SchoolName).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("学校名称已存在")
	}

	// 如果指定了负责人，检查用户是否存在
	if req.SchoolAdmin != nil {
		var user models.User
		if err := config.DB.First(&user, *req.SchoolAdmin).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("指定的负责人不存在")
			}
			return nil, err
		}
	}

	// 创建学校
	school := &models.School{
		SchoolName:  req.SchoolName,
		SchoolAdmin: req.SchoolAdmin,
		UserCount:   0,
		Status:      "active",
	}

	if err := config.DB.Create(school).Error; err != nil {
		return nil, err
	}

	return school, nil
}

// GetSchoolList 获取学校列表
func (s *SchoolService) GetSchoolList(req *dto.SchoolListRequest) (*dto.SchoolListResponse, error) {
	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Order == "" {
		req.Order = "desc"
	}

	// 构建查询
	query := config.DB.Model(&models.School{})

	// 搜索条件
	if req.Search != "" {
		query = query.Where("school_name LIKE ?", "%"+req.Search+"%")
	}

	// 状态筛选
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 排序
	sortBy := "created_at"
	if req.SortBy != "" {
		sortBy = req.SortBy
	}
	orderSQL := sortBy + " " + strings.ToUpper(req.Order)
	query = query.Order(orderSQL)

	// 分页查询
	offset := (req.Page - 1) * req.Limit
	var schools []models.School
	if err := query.Offset(offset).Limit(req.Limit).Find(&schools).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]dto.SchoolResponse, 0, len(schools))
	for _, school := range schools {
		resp := dto.SchoolResponse{
			ID:          school.ID,
			SchoolName:  school.SchoolName,
			SchoolAdmin: school.SchoolAdmin,
			UserCount:   school.UserCount,
			Status:      school.Status,
			CreatedAt:   school.CreatedAt,
			UpdatedAt:   school.UpdatedAt,
		}

		// 查询负责人姓名
		if school.SchoolAdmin != nil {
			var user models.User
			if err := config.DB.Select("username").First(&user, *school.SchoolAdmin).Error; err == nil {
				resp.AdminName = &user.Username
			}
		}

		list = append(list, resp)
	}

	return &dto.SchoolListResponse{
		Total: int(total),
		Page:  req.Page,
		Limit: req.Limit,
		List:  list,
	}, nil
}

// GetSchoolByID 根据ID获取学校
func (s *SchoolService) GetSchoolByID(id int64) (*dto.SchoolResponse, error) {
	var school models.School
	if err := config.DB.First(&school, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("学校不存在")
		}
		return nil, err
	}

	resp := &dto.SchoolResponse{
		ID:          school.ID,
		SchoolName:  school.SchoolName,
		SchoolAdmin: school.SchoolAdmin,
		UserCount:   school.UserCount,
		Status:      school.Status,
		CreatedAt:   school.CreatedAt,
		UpdatedAt:   school.UpdatedAt,
	}

	// 查询负责人姓名
	if school.SchoolAdmin != nil {
		var user models.User
		if err := config.DB.Select("username").First(&user, *school.SchoolAdmin).Error; err == nil {
			resp.AdminName = &user.Username
		}
	}

	return resp, nil
}

// GetSchoolByName 根据名称获取学校
func (s *SchoolService) GetSchoolByName(name string) (*dto.SchoolResponse, error) {
	var school models.School
	if err := config.DB.Where("school_name = ?", name).First(&school).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("学校不存在")
		}
		return nil, err
	}

	resp := &dto.SchoolResponse{
		ID:          school.ID,
		SchoolName:  school.SchoolName,
		SchoolAdmin: school.SchoolAdmin,
		UserCount:   school.UserCount,
		Status:      school.Status,
		CreatedAt:   school.CreatedAt,
		UpdatedAt:   school.UpdatedAt,
	}

	// 查询负责人姓名
	if school.SchoolAdmin != nil {
		var user models.User
		if err := config.DB.Select("username").First(&user, *school.SchoolAdmin).Error; err == nil {
			resp.AdminName = &user.Username
		}
	}

	return resp, nil
}

// UpdateSchool 更新学校信息
func (s *SchoolService) UpdateSchool(id int64, req *dto.UpdateSchoolRequest) (*models.School, error) {
	var school models.School
	if err := config.DB.First(&school, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("学校不存在")
		}
		return nil, err
	}

	// 更新字段
	updates := make(map[string]interface{})

	if req.SchoolName != nil {
		// 检查新名称是否已被使用
		var count int64
		if err := config.DB.Model(&models.School{}).
			Where("school_name = ? AND id != ?", *req.SchoolName, id).
			Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("学校名称已存在")
		}
		updates["school_name"] = *req.SchoolName
	}

	if req.SchoolAdmin != nil {
		// 检查用户是否存在
		var user models.User
		if err := config.DB.First(&user, *req.SchoolAdmin).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("指定的负责人不存在")
			}
			return nil, err
		}
		updates["school_admin"] = *req.SchoolAdmin
	}

	if len(updates) > 0 {
		if err := config.DB.Model(&school).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// 重新查询获取最新数据
	if err := config.DB.First(&school, id).Error; err != nil {
		return nil, err
	}

	return &school, nil
}

// UpdateSchoolStatus 更新学校状态
func (s *SchoolService) UpdateSchoolStatus(id int64, status string) error {
	var school models.School
	if err := config.DB.First(&school, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("学校不存在")
		}
		return err
	}

	if err := config.DB.Model(&school).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

// DeleteSchool 删除学校（软删除）
func (s *SchoolService) DeleteSchool(id int64) error {
	var school models.School
	if err := config.DB.First(&school, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("学校不存在")
		}
		return err
	}

	// 检查是否有学生关联
	var userCount int64
	if err := config.DB.Model(&models.User{}).Where("school_id = ?", id).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		return errors.New("该学校还有学生关联，无法删除")
	}

	// 软删除
	if err := config.DB.Delete(&school).Error; err != nil {
		return err
	}

	return nil
}
