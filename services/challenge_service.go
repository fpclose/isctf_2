package services

import (
	"errors"
	"fmt"
	"isctf/config"
	"isctf/dto"
	"isctf/models"
	"isctf/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ChallengeService struct{}

func NewChallengeService() *ChallengeService {
	return &ChallengeService{}
}

// CreateChallenge 创建题目
func (s *ChallengeService) CreateChallenge(req *dto.ChallengeRequest) (*models.Challenge, error) {
	chal := &models.Challenge{
		ChallengeName: req.ChallengeName,
		Direction:     req.Direction,
		Author:        req.Author,
		Description:   req.Description,
		Hint:          req.Hint,
		State:         req.State,
		Mode:          req.Mode,
		StaticFlag:    req.StaticFlag,
		DockerImage:   req.DockerImage,
		Difficulty:    req.Difficulty,
		InitialScore:  req.InitialScore,
		MinScore:      req.MinScore,
		CurrentScore:  req.InitialScore,
		DecayRatio:    req.DecayRatio,
		SolvedCount:   0,
	}
	if req.DockerPorts != nil {
		chal.DockerPorts = models.DockerPorts(req.DockerPorts)
	}

	if err := config.DB.Create(chal).Error; err != nil {
		return nil, err
	}
	return chal, nil
}

// GetChallengeList 获取列表
func (s *ChallengeService) GetChallengeList(req *dto.ChallengeListRequest, isAdmin bool) (interface{}, int64, error) {
	var list []models.Challenge
	var total int64
	db := config.DB.Model(&models.Challenge{})

	if !isAdmin {
		db = db.Where("state = ?", "visible")
	} else if req.State != "" {
		db = db.Where("state = ?", req.State)
	}

	if req.Direction != "" {
		db = db.Where("direction = ?", req.Direction)
	}
	if req.Difficulty != "" {
		db = db.Where("difficulty = ?", req.Difficulty)
	}
	if req.Search != "" {
		db = db.Where("challenge_name LIKE ?", "%"+req.Search+"%")
	}

	db.Count(&total)
	err := db.Order("id DESC").Offset((req.Page - 1) * req.Limit).Limit(req.Limit).Find(&list).Error
	return list, total, err
}

// GetDetail 获取详情
func (s *ChallengeService) GetDetail(id int64, isAdmin bool) (*models.Challenge, error) {
	var chal models.Challenge
	db := config.DB
	if !isAdmin {
		db = db.Where("state = ?", "visible")
	}
	if err := db.First(&chal, id).Error; err != nil {
		return nil, errors.New("题目不存在或不可见")
	}
	return &chal, nil
}

// StartContainer 启动动态题目容器
func (s *ChallengeService) StartContainer(userID, teamID, challengeID int64) (*models.Container, error) {
	// 1. 检查题目
	var chal models.Challenge
	if err := config.DB.First(&chal, challengeID).Error; err != nil {
		return nil, errors.New("题目不存在")
	}
	if chal.Mode != "dynamic" {
		return nil, errors.New("非动态题目无需启动容器")
	}
	if chal.State != "visible" {
		return nil, errors.New("题目不可见")
	}

	// 2. 检查是否已有运行中的容器
	var existing models.Container
	err := config.DB.Where("team_id = ? AND challenge_id = ? AND state = 'running'", teamID, challengeID).First(&existing).Error
	if err == nil {
		// 已存在
		if time.Now().Before(existing.EndTime) {
			return &existing, nil
		}
		// 已过期但状态未更新，先停止
		_ = s.stopContainerInternal(&existing)
	}

	// 3. 生成 Flag
	flag := utils.GenerateDynamicFlag(teamID, challengeID)

	// 4. 调用 Docker API
	// 环境变量注入 Flag
	env := []string{
		fmt.Sprintf("FLAG=%s", flag),
		fmt.Sprintf("GZCTF_FLAG=%s", flag), // 兼容常见CTF镜像
	}

	containerID, hostMapping, err := utils.StartContainer(*chal.DockerImage, chal.DockerPorts, env)
	if err != nil {
		return nil, fmt.Errorf("启动容器失败: %v", err)
	}

	// 5. 记录数据库
	// 注意：Docker 返回的 containerID 是长 ID
	newContainer := &models.Container{
		ChallengeID:   challengeID,
		TeamID:        teamID,
		UserID:        userID,
		ContainerName: containerID, // 存储 Docker ID
		DockerImage:   *chal.DockerImage,
		DockerPorts:   chal.DockerPorts,
		HostMapping:   models.PortMapping(hostMapping),
		ContainerFlag: flag,
		State:         "running",
		StartTime:     time.Now(),
		EndTime:       time.Now().Add(1 * time.Hour), // 默认1小时
		ExtendedCount: 0,
	}

	if err := config.DB.Create(newContainer).Error; err != nil {
		// 数据库插入失败，回滚Docker操作
		_ = utils.RemoveContainer(containerID)
		return nil, err
	}

	return newContainer, nil
}

// StopChallengeContainer 停止题目容器
func (s *ChallengeService) StopChallengeContainer(userID, challengeID int64) error {
	// 查找该用户(所属团队)在该题目的运行容器
	// 这里简化为查找用户启动的，实际上应该查找团队的
	// 需要 TeamService 辅助查找 TeamID，这里简化处理
	var c models.Container
	// 假设能找到归属
	if err := config.DB.Where("user_id = ? AND challenge_id = ? AND state = 'running'", userID, challengeID).First(&c).Error; err != nil {
		return errors.New("未找到运行中的容器")
	}
	return s.stopContainerInternal(&c)
}

func (s *ChallengeService) stopContainerInternal(c *models.Container) error {
	// 调用 Docker API
	utils.RemoveContainer(c.ContainerName)

	// 更新数据库
	c.State = "destroyed"
	return config.DB.Save(c).Error
}

// SubmitFlag 提交 Flag
func (s *ChallengeService) SubmitFlag(userID, teamID, challengeID int64, flag string, ip string) (bool, int, error) {
	// 1. 获取题目
	var chal models.Challenge
	if err := config.DB.First(&chal, challengeID).Error; err != nil {
		return false, 0, errors.New("题目不存在")
	}

	// 2. 检查是否已解出
	var solveCount int64
	config.DB.Model(&models.Solve{}).Where("team_id = ? AND challenge_id = ?", teamID, challengeID).Count(&solveCount)
	if solveCount > 0 {
		return false, 0, errors.New("本团队已解出该题")
	}

	// 3. 验证 Flag
	var isCorrect bool
	flag = strings.TrimSpace(flag)

	if chal.Mode == "static" {
		isCorrect = (chal.StaticFlag != nil && *chal.StaticFlag == flag)
	} else {
		// 动态题，查找该团队的容器 Flag
		var container models.Container
		err := config.DB.Where("team_id = ? AND challenge_id = ? AND state = 'running'", teamID, challengeID).First(&container).Error
		if err != nil {
			// 如果容器已销毁，是否允许提交？通常不允许，或者查最近记录。
			// 这里简单处理：允许查最近一个非 destroyed 状态的（防止刚好过期）或最近的一条记录
			err = config.DB.Where("team_id = ? AND challenge_id = ?", teamID, challengeID).Order("id desc").First(&container).Error
		}

		if err == nil {
			isCorrect = (container.ContainerFlag == flag)
		} else {
			isCorrect = false
		}
	}

	// 4. 记录日志
	log := &models.SubmissionLog{
		ChallengeID:   challengeID,
		TeamID:        teamID,
		UserID:        userID,
		SubmittedFlag: flag,
		FlagResult:    "wrong",
		ChallengeType: chal.Mode,
		IPAddress:     ip,
	}
	if isCorrect {
		log.FlagResult = "correct"
	}
	config.DB.Create(log)

	if !isCorrect {
		return false, 0, nil
	}

	// 5. 处理解出逻辑（事务）
	score := 0
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// 锁定题目行以更新分数
		var lockedChal models.Challenge
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&lockedChal, challengeID).Error; err != nil {
			return err
		}

		score = lockedChal.CurrentScore

		// 记录 Solve
		solve := &models.Solve{
			ChallengeID:   challengeID,
			TeamID:        teamID,
			UserID:        userID,
			EarnedScore:   score,
			Rank:          lockedChal.SolvedCount + 1,
			IsFirstBlood:  lockedChal.SolvedCount == 0,
			IsSecondBlood: lockedChal.SolvedCount == 1,
			IsThirdBlood:  lockedChal.SolvedCount == 2,
		}
		if err := tx.Create(solve).Error; err != nil {
			return err
		}

		// 更新题目状态（解出数+1，分数衰减）
		if err := lockedChal.UpdateScore(tx); err != nil {
			return err
		}

		// 更新团队总分
		if err := tx.Model(&models.Team{}).Where("id = ?", teamID).
			UpdateColumn("team_score", gorm.Expr("team_score + ?", score)).Error; err != nil {
			return err
		}

		return nil
	})

	return true, score, err
}
