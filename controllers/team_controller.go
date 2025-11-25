package controllers

import (
	"isctf/dto"
	"isctf/services"
	"isctf/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TeamController 团队控制器
type TeamController struct {
	teamService *services.TeamService
}

// NewTeamController 创建团队控制器实例
func NewTeamController() *TeamController {
	return &TeamController{
		teamService: services.NewTeamService(),
	}
}

// CreateTeam 创建团队
func (c *TeamController) CreateTeam(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Error(ctx, utils.UNAUTHORIZED)
		return
	}

	var req dto.CreateTeamRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	team, err := c.teamService.CreateTeam(userID.(int64), &req)
	if err != nil {
		if err.Error() == "团队名称已存在" {
			utils.ErrorWithMsg(ctx, utils.TEAM_ALREADY_EXIST, err.Error())
			return
		}
		if err.Error() == "您已经加入了其他团队" {
			utils.ErrorWithMsg(ctx, utils.TEAM_ALREADY_JOINED, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "创建团队失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "创建团队成功", gin.H{
		"id":           team.ID,
		"team_name":    team.TeamName,
		"captain_name": team.CaptainName,
		"team_track":   team.TeamTrack,
		"member_count": team.MemberCount,
	})
}

// JoinTeam 加入团队
func (c *TeamController) JoinTeam(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Error(ctx, utils.UNAUTHORIZED)
		return
	}

	var req dto.JoinTeamRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	team, err := c.teamService.JoinTeam(userID.(int64), &req)
	if err != nil {
		if err.Error() == "团队不存在" {
			utils.ErrorWithMsg(ctx, utils.TEAM_NOT_EXIST, err.Error())
			return
		}
		if err.Error() == "团队密码错误" {
			utils.ErrorWithMsg(ctx, utils.TEAM_PASSWORD_ERROR, err.Error())
			return
		}
		if err.Error() == "团队已满员" {
			utils.ErrorWithMsg(ctx, utils.TEAM_FULL, err.Error())
			return
		}
		if err.Error() == "您已经加入了其他团队" {
			utils.ErrorWithMsg(ctx, utils.TEAM_ALREADY_JOINED, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "加入团队失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "加入团队成功", gin.H{
		"id":           team.ID,
		"team_name":    team.TeamName,
		"captain_name": team.CaptainName,
		"member_count": team.MemberCount,
	})
}

// LeaveTeam 离开团队
func (c *TeamController) LeaveTeam(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Error(ctx, utils.UNAUTHORIZED)
		return
	}

	if err := c.teamService.LeaveTeam(userID.(int64)); err != nil {
		if err.Error() == "您还未加入任何团队" {
			utils.ErrorWithMsg(ctx, utils.TEAM_NOT_JOINED, err.Error())
			return
		}
		if err.Error() == "队长无法直接退出，请先转让队长或解散团队" {
			utils.ErrorWithMsg(ctx, utils.TEAM_CAPTAIN_CANNOT_LEAVE, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "离开团队失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "离开团队成功", nil)
}

// GetMyTeam 获取我的团队信息
func (c *TeamController) GetMyTeam(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Error(ctx, utils.UNAUTHORIZED)
		return
	}

	team, err := c.teamService.GetMyTeam(userID.(int64))
	if err != nil {
		if err.Error() == "您还未加入任何团队" {
			utils.ErrorWithMsg(ctx, utils.TEAM_NOT_JOINED, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "获取团队信息失败: "+err.Error())
		return
	}

	utils.Success(ctx, team)
}

// GetTeamByID 根据ID获取团队信息
func (c *TeamController) GetTeamByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	teamID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || teamID <= 0 {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "无效的团队ID")
		return
	}

	team, err := c.teamService.GetTeamByID(teamID)
	if err != nil {
		if err.Error() == "团队不存在" {
			utils.Error(ctx, utils.TEAM_NOT_EXIST)
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "获取团队信息失败")
		return
	}

	utils.Success(ctx, team)
}

// GetTeamList 获取团队列表
func (c *TeamController) GetTeamList(ctx *gin.Context) {
	var req dto.TeamListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	result, err := c.teamService.GetTeamList(&req)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.ERROR, "获取团队列表失败: "+err.Error())
		return
	}

	utils.Success(ctx, result)
}

// UpdateTeam 更新团队信息（队长）
func (c *TeamController) UpdateTeam(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Error(ctx, utils.UNAUTHORIZED)
		return
	}

	var req dto.UpdateTeamRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	if err := c.teamService.UpdateTeam(userID.(int64), &req); err != nil {
		if err.Error() == "您不是任何团队的队长" {
			utils.ErrorWithMsg(ctx, utils.TEAM_NOT_CAPTAIN, err.Error())
			return
		}
		if err.Error() == "团队名称已存在" {
			utils.ErrorWithMsg(ctx, utils.TEAM_ALREADY_EXIST, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "更新团队失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "更新团队成功", nil)
}

// TransferCaptain 转让队长（队长）
func (c *TeamController) TransferCaptain(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Error(ctx, utils.UNAUTHORIZED)
		return
	}

	var req dto.TransferCaptainRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	if err := c.teamService.TransferCaptain(userID.(int64), &req); err != nil {
		if err.Error() == "您不是任何团队的队长" {
			utils.ErrorWithMsg(ctx, utils.TEAM_NOT_CAPTAIN, err.Error())
			return
		}
		if err.Error() == "新队长必须是团队成员" {
			utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "转让队长失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "转让队长成功", nil)
}

// RemoveMember 移除成员（队长）
func (c *TeamController) RemoveMember(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Error(ctx, utils.UNAUTHORIZED)
		return
	}

	var req dto.RemoveMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	if err := c.teamService.RemoveMember(userID.(int64), &req); err != nil {
		if err.Error() == "您不是任何团队的队长" {
			utils.ErrorWithMsg(ctx, utils.TEAM_NOT_CAPTAIN, err.Error())
			return
		}
		if err.Error() == "该用户不是团队成员" {
			utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "移除成员失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "移除成员成功", nil)
}

// DisbandTeam 解散团队（队长）
func (c *TeamController) DisbandTeam(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Error(ctx, utils.UNAUTHORIZED)
		return
	}

	if err := c.teamService.DisbandTeam(userID.(int64)); err != nil {
		if err.Error() == "您不是任何团队的队长" {
			utils.ErrorWithMsg(ctx, utils.TEAM_NOT_CAPTAIN, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "解散团队失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "解散团队成功", nil)
}

// UpdateTeamStatus 更新团队状态（管理员）
func (c *TeamController) UpdateTeamStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	teamID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "无效的团队ID")
		return
	}

	var req dto.UpdateTeamStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	if err := c.teamService.UpdateTeamStatus(teamID, req.Status); err != nil {
		if err.Error() == "团队不存在" {
			utils.ErrorWithMsg(ctx, utils.TEAM_NOT_EXIST, err.Error())
			return
		}
		utils.ErrorWithMsg(ctx, utils.ERROR, "更新团队状态失败: "+err.Error())
		return
	}

	utils.SuccessWithMsg(ctx, "更新团队状态成功", nil)
}

// GetTeamRank 获取团队排名
func (c *TeamController) GetTeamRank(ctx *gin.Context) {
	var req dto.TeamRankRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, "参数错误: "+err.Error())
		return
	}

	result, err := c.teamService.GetTeamRank(&req)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.ERROR, "获取团队排名失败: "+err.Error())
		return
	}

	utils.Success(ctx, result)
}
