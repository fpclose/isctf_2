package controllers

import (
	"isctf/config"
	"isctf/dto"
	"isctf/services"
	"isctf/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChallengeController struct {
	chalService *services.ChallengeService
}

func NewChallengeController() *ChallengeController {
	return &ChallengeController{
		chalService: services.NewChallengeService(),
	}
}

// Create 创建题目
func (c *ChallengeController) Create(ctx *gin.Context) {
	var req dto.ChallengeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, err.Error())
		return
	}
	chal, err := c.chalService.CreateChallenge(&req)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.ERROR, err.Error())
		return
	}
	utils.Success(ctx, chal)
}

// Update 更新题目
func (c *ChallengeController) Update(ctx *gin.Context) {
	// 略: 实现更新逻辑
	utils.Success(ctx, nil)
}

// Delete 删除题目
func (c *ChallengeController) Delete(ctx *gin.Context) {
	// 略: 实现删除逻辑
	utils.Success(ctx, nil)
}

// UpdateState 更新状态
func (c *ChallengeController) UpdateState(ctx *gin.Context) {
	// 略: 实现更新状态逻辑
	utils.Success(ctx, nil)
}

// GetList 获取题目列表
func (c *ChallengeController) GetList(ctx *gin.Context) {
	var req dto.ChallengeListRequest
	ctx.ShouldBindQuery(&req)

	role, _ := ctx.Get("role")
	isAdmin := (role == "admin" || role == "super_admin")

	list, total, err := c.chalService.GetChallengeList(&req, isAdmin)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.ERROR, err.Error())
		return
	}
	utils.Success(ctx, gin.H{"list": list, "total": total})
}

// GetDetail 获取详情
func (c *ChallengeController) GetDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	role, _ := ctx.Get("role")
	isAdmin := (role == "admin" || role == "super_admin")

	chal, err := c.chalService.GetDetail(id, isAdmin)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.ERROR, err.Error())
		return
	}
	utils.Success(ctx, chal)
}

// StartContainer 启动容器
func (c *ChallengeController) StartContainer(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	// 必须已加入团队
	ts := services.NewTeamService()
	teamDetail, err := ts.GetMyTeam(userID)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.TEAM_NOT_JOINED, "未加入团队无法启动环境")
		return
	}

	idStr := ctx.Param("id")
	chalID, _ := strconv.ParseInt(idStr, 10, 64)

	container, err := c.chalService.StartContainer(userID, teamDetail.ID, chalID)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.ERROR, err.Error())
		return
	}

	// 构造前端友好的返回
	hostIP := config.AppConfig.Database.Host // 简化：生产环境应配置 Docker 主机的公网 IP
	if hostIP == "127.0.0.1" || hostIP == "localhost" {
		hostIP = ctx.Request.Host // 开发环境尝试取请求 Host
	}

	utils.Success(ctx, gin.H{
		"container_id":   container.ID,
		"status":         "running",
		"host":           hostIP,
		"ports":          container.HostMapping,
		"expires_at":     container.EndTime,
		"time_remaining": "1h", // 简化
	})
}

// StopContainer 停止容器
func (c *ChallengeController) StopContainer(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	idStr := ctx.Param("id")
	chalID, _ := strconv.ParseInt(idStr, 10, 64)

	if err := c.chalService.StopChallengeContainer(userID, chalID); err != nil {
		utils.ErrorWithMsg(ctx, utils.ERROR, err.Error())
		return
	}
	utils.SuccessWithMsg(ctx, "环境已停止", nil)
}

// RenewContainer 续期容器
func (c *ChallengeController) RenewContainer(ctx *gin.Context) {
	// 略
	utils.Success(ctx, nil)
}

// GetContainerStatus 获取状态
func (c *ChallengeController) GetContainerStatus(ctx *gin.Context) {
	// 略
	utils.Success(ctx, nil)
}

// DestroyContainer 销毁容器 (通用接口)
func (c *ChallengeController) DestroyContainer(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	// 这里的 id 可能是容器 ID，也可能是题目 ID，需根据 API 定义调整
	// 这里假设 POST body 传 container_id 或 challenge_id
	// 暂时简单实现
	utils.SuccessWithMsg(ctx, "容器已销毁", nil)
}

// SubmitFlag 提交 Flag
func (c *ChallengeController) SubmitFlag(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	idStr := ctx.Param("id")
	chalID, _ := strconv.ParseInt(idStr, 10, 64)

	var req dto.SubmitFlagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithMsg(ctx, utils.INVALID_PARAMS, err.Error())
		return
	}

	ts := services.NewTeamService()
	teamDetail, err := ts.GetMyTeam(userID)
	if err != nil {
		utils.ErrorWithMsg(ctx, utils.TEAM_NOT_JOINED, "未加入团队无法提交")
		return
	}

	correct, score, err := c.chalService.SubmitFlag(userID, teamDetail.ID, chalID, req.Flag, ctx.ClientIP())
	if err != nil {
		// 可能是已解答或题目不存在
		utils.ErrorWithMsg(ctx, utils.ERROR, err.Error())
		return
	}

	if correct {
		utils.SuccessWithMsg(ctx, "Flag 正确！", gin.H{"earned_score": score})
	} else {
		utils.ErrorWithMsg(ctx, 432, "Flag 错误") // 432 自定义错误码
	}
}

// Attachment 相关接口略... (UploadAttachment, DeleteAttachment, DownloadAttachment)
func (c *ChallengeController) UploadAttachment(ctx *gin.Context)   {}
func (c *ChallengeController) DeleteAttachment(ctx *gin.Context)   {}
func (c *ChallengeController) DownloadAttachment(ctx *gin.Context) {}

// Admin 相关接口略...
func (c *ChallengeController) GetAdminContainers(ctx *gin.Context) {}
func (c *ChallengeController) AdminStopContainer(ctx *gin.Context) {}

// Log 相关接口略...
func (c *ChallengeController) GetRecentSolves(ctx *gin.Context) {}
func (c *ChallengeController) GetTeamSolves(ctx *gin.Context)   {}
func (c *ChallengeController) GetUserSolves(ctx *gin.Context)   {}
