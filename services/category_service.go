package services

import (
	"errors"
	"isctf/config"
	"isctf/dto"
	"isctf/models"
)

type CategoryService struct{}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

func (s *CategoryService) CreateCategory(req *dto.CategoryRequest) (*models.ChallengeCategory, error) {
	cat := &models.ChallengeCategory{
		Direction:   req.Direction,
		NameZh:      req.NameZh,
		NameEn:      req.NameEn,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
		SortOrder:   req.SortOrder,
		Status:      "active",
	}
	if err := config.DB.Create(cat).Error; err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) GetList() ([]models.ChallengeCategory, error) {
	var list []models.ChallengeCategory
	// 按 sort_order 排序
	err := config.DB.Order("sort_order ASC").Find(&list).Error
	return list, err
}

// 其他 CRUD 省略...
