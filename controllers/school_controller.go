package controllers

import (
	"isctf/dto"
	"isctf/services"
	"isctf/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SchoolController 学校控制器
type SchoolController struct {
	schoolService *services.SchoolService
}

// NewSchoolController 创建学校控制器实例
func NewSchoolController() *SchoolController {
	return &SchoolController{
		schoolService: services.NewSchoolService(),
	}
}

// CreateSchool 创建学校
// @Summary 创建学校
// @Tags 学校管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param request body dto.CreateSchoolRequest true "学校信息"
// @Success 200 {object} utils.Response
// @Router /api/v1/schools [post]
func (c *SchoolController) CreateSchool(ctx *gin.Context) {
	var req dto.CreateSchoolRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	school, err := c.schoolService.CreateSchool(&req)
	if err != nil {
		if err.Error() == "学校名称已存在" {
			utils.ErrorWithMsg(ctx, utils.SCHOOL_ALREADY_EXIST, err.Error())
			return
		}
		if err.Error() == "指定的负责人不存在" {
			utils.ErrorWithMsg(ctx, utils.USER_NOT_EXIST, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "创建学校失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "创建学校成功", gin.H{
		"id":          school.ID,
		"school_name": school.SchoolName,
		"user_count":  school.UserCount,
		"status":      school.Status,
	})
}

// GetSchools 获取学校列表
// @Summary 获取学校列表
// @Tags 学校管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param limit query int false "每页数量"
// @Param search query string false "搜索关键词"
// @Param sort_by query string false "排序字段"
// @Param order query string false "排序方式"
// @Param status query string false "状态筛选"
// @Success 200 {object} utils.Response
// @Router /api/v1/schools [get]
func (c *SchoolController) GetSchools(ctx *gin.Context) {
	var req dto.SchoolListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	result, err := c.schoolService.GetSchoolList(&req)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.ERROR, "获取学校列表失败: "+err.Error())
		return
	}

	utils.Success(ctx, result)
}

// GetSchoolByID 根据ID获取学校详情
// @Summary 获取学校详情
// @Tags 学校管理
// @Accept json
// @Produce json
// @Param id path int true "学校ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/schools/:id [get]
func (c *SchoolController) GetSchoolByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "无效的学校ID")
		return
	}

	school, err := c.schoolService.GetSchoolByID(id)
	if err != nil {
		if err.Error() == "学校不存在" {
			utils.ErrorWithMsg(ctx, utils.SCHOOL_NOT_EXIST, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "获取学校详情失败: "+err.Error())
		return
	}

	utils.Success(ctx, school)
}

// GetSchoolByName 根据名称获取学校详情
// @Summary 根据名称获取学校详情
// @Tags 学校管理
// @Accept json
// @Produce json
// @Param school_name path string true "学校名称"
// @Success 200 {object} utils.Response
// @Router /api/v1/schools/name/:school_name [get]
func (c *SchoolController) GetSchoolByName(ctx *gin.Context) {
	schoolName := ctx.Param("school_name")
	if schoolName == "" {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "学校名称不能为空")
		return
	}

	school, err := c.schoolService.GetSchoolByName(schoolName)
	if err != nil {
		if err.Error() == "学校不存在" {
			utils.ErrorWithMsg(ctx, utils.SCHOOL_NOT_EXIST, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "获取学校详情失败: "+err.Error())
		return
	}

	utils.Success(ctx, school)
}

// UpdateSchool 更新学校信息
// @Summary 更新学校信息
// @Tags 学校管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "学校ID"
// @Param request body dto.UpdateSchoolRequest true "学校信息"
// @Success 200 {object} utils.Response
// @Router /api/v1/schools/:id [put]
func (c *SchoolController) UpdateSchool(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "无效的学校ID")
		return
	}

	var req dto.UpdateSchoolRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	school, err := c.schoolService.UpdateSchool(id, &req)
	if err != nil {
		if err.Error() == "学校不存在" {
			utils.ErrorWithMsg(ctx, utils.SCHOOL_NOT_EXIST, err.Error())
			return
		}
		if err.Error() == "学校名称已存在" {
			utils.ErrorWithMsg(ctx, utils.SCHOOL_ALREADY_EXIST, err.Error())
			return
		}
		if err.Error() == "指定的负责人不存在" {
			utils.ErrorWithMsg(ctx, utils.USER_NOT_EXIST, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "更新学校失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "更新学校成功", gin.H{
		"id":          school.ID,
		"school_name": school.SchoolName,
		"user_count":  school.UserCount,
		"status":      school.Status,
	})
}

// UpdateSchoolStatus 更新学校状态
// @Summary 更新学校状态
// @Tags 学校管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "学校ID"
// @Param request body dto.UpdateSchoolStatusRequest true "状态信息"
// @Success 200 {object} utils.Response
// @Router /api/v1/schools/:id/status [patch]
func (c *SchoolController) UpdateSchoolStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "无效的学校ID")
		return
	}

	var req dto.UpdateSchoolStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	if err := c.schoolService.UpdateSchoolStatus(id, req.Status); err != nil {
		if err.Error() == "学校不存在" {
			utils.ErrorWithMsg(ctx, utils.SCHOOL_NOT_EXIST, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "更新学校状态失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "更新学校状态成功", nil)
}

// DeleteSchool 删除学校
// @Summary 删除学校
// @Tags 学校管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "学校ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/schools/:id [delete]
func (c *SchoolController) DeleteSchool(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "无效的学校ID")
		return
	}

	if err := c.schoolService.DeleteSchool(id); err != nil {
		if err.Error() == "学校不存在" {
			utils.ErrorWithMsg(ctx, utils.SCHOOL_NOT_EXIST, err.Error())
			return
		}
		if err.Error() == "该学校还有学生关联，无法删除" {
			utils.ErrorWithMsg(ctx, utils.CONFLICT, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "删除学校失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "删除学校成功", nil)
}
