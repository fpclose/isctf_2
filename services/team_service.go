package services

import (
	"errors"
	"isctf/config"
	"isctf/dto"
	"isctf/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

// TeamService 团队服务
type TeamService struct{}

// NewTeamService 创建团队服务实例
func NewTeamService() *TeamService {
	return &TeamService{}
}

// CreateTeam 创建团队
func (s *TeamService) CreateTeam(userID int64, req *dto.CreateTeamRequest) (*models.Team, error) {
	// 检查团队名称是否已存在
	var count int64
	if err := config.DB.Model(&models.Team{}).Where("team_name = ?", req.TeamName).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("团队名称已存在")
	}

	// 获取用户信息
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户是否已经在其他团队
	if err := config.DB.Model(&models.Team{}).
		Where("(captain_id = ? OR member1_id = ? OR member2_id = ?) AND status = 'active'", userID, userID, userID).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("您已经加入了其他团队")
	}

	// 验证赛道匹配
	if req.TeamTrack == "social" && user.Track != "social" {
		return nil, errors.New("只有社会赛道用户才能创建社会赛道团队")
	}
	if (req.TeamTrack == "freshman" || req.TeamTrack == "advanced") && user.Track != "school" {
		return nil, errors.New("只有联合院校赛道用户才能创建校内团队")
	}

	// 创建团队
	team := &models.Team{
		TeamName:    req.TeamName,
		CaptainID:   userID,
		CaptainName: user.Username,
		TeamTrack:   req.TeamTrack,
		SchoolID:    user.SchoolID,
		SchoolName:  user.SchoolName,
		TeamScore:   0,
		MemberCount: 1,
		Status:      "active",
	}

	// 设置密码
	if err := team.SetPassword(req.TeamPassword); err != nil {
		return nil, err
	}

	if err := config.DB.Create(team).Error; err != nil {
		return nil, err
	}

	return team, nil
}

// JoinTeam 加入团队
func (s *TeamService) JoinTeam(userID int64, req *dto.JoinTeamRequest) (*models.Team, error) {
	// 查找团队
	var team models.Team
	if err := config.DB.Where("team_name = ?", req.TeamName).First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("团队不存在")
		}
		return nil, err
	}

	// 检查团队状态
	if team.Status != "active" {
		return nil, errors.New("该团队已解散或被封禁，无法加入")
	}

	// 验证密码
	if !team.CheckPassword(req.TeamPassword) {
		return nil, errors.New("团队密码错误")
	}

	// 检查团队是否已满员
	if team.IsFull() {
		return nil, errors.New("团队已满员")
	}

	// 获取用户信息
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户是否已经在其他团队
	var count int64
	if err := config.DB.Model(&models.Team{}).
		Where("(captain_id = ? OR member1_id = ? OR member2_id = ?) AND status = 'active'", userID, userID, userID).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("您已经加入了其他团队")
	}

	// 检查赛道是否匹配
	if team.TeamTrack == "social" && user.Track != "social" {
		return nil, errors.New("赛道不匹配，无法加入该团队")
	}
	if (team.TeamTrack == "freshman" || team.TeamTrack == "advanced") && user.Track != "school" {
		return nil, errors.New("赛道不匹配，无法加入该团队")
	}

	// 检查学校是否匹配（校内赛道）
	if team.TeamTrack == "freshman" || team.TeamTrack == "advanced" {
		if team.SchoolID == nil || user.SchoolID == nil || *team.SchoolID != *user.SchoolID {
			return nil, errors.New("校内赛道团队只能加入同一学校的成员")
		}
	}

	// 添加成员
	updates := make(map[string]interface{})
	if team.Member1ID == nil {
		updates["member1_id"] = userID
		updates["member1_name"] = user.Username
	} else if team.Member2ID == nil {
		updates["member2_id"] = userID
		updates["member2_name"] = user.Username
	}
	updates["member_count"] = team.MemberCount + 1

	if err := config.DB.Model(&team).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 重新查询获取最新数据
	if err := config.DB.First(&team, team.ID).Error; err != nil {
		return nil, err
	}

	return &team, nil
}

// LeaveTeam 离开团队
func (s *TeamService) LeaveTeam(userID int64) error {
	// 查找用户所在的团队
	var team models.Team
	if err := config.DB.Where("(captain_id = ? OR member1_id = ? OR member2_id = ?) AND status = 'active'",
		userID, userID, userID).First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("您还未加入任何团队")
		}
		return err
	}

	// 队长不能直接退出，需要先转让队长或解散团队
	if team.CaptainID == userID {
		return errors.New("队长无法直接退出，请先转让队长或解散团队")
	}

	// 移除成员
	updates := make(map[string]interface{})
	if team.Member1ID != nil && *team.Member1ID == userID {
		updates["member1_id"] = nil
		updates["member1_name"] = nil
	} else if team.Member2ID != nil && *team.Member2ID == userID {
		updates["member2_id"] = nil
		updates["member2_name"] = nil
	}
	updates["member_count"] = team.MemberCount - 1

	if err := config.DB.Model(&team).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// GetMyTeam 获取我的团队信息
func (s *TeamService) GetMyTeam(userID int64) (*dto.TeamDetailResponse, error) {
	// 查找用户所在的团队
	var team models.Team
	if err := config.DB.Where("(captain_id = ? OR member1_id = ? OR member2_id = ?) AND status = 'active'",
		userID, userID, userID).First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("您还未加入任何团队")
		}
		return nil, err
	}

	// 构建团队详情
	return s.buildTeamDetail(&team)
}

// GetTeamByID 根据ID获取团队信息
func (s *TeamService) GetTeamByID(teamID int64) (*dto.TeamDetailResponse, error) {
	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("团队不存在")
		}
		return nil, err
	}

	return s.buildTeamDetail(&team)
}

// buildTeamDetail 构建团队详情
func (s *TeamService) buildTeamDetail(team *models.Team) (*dto.TeamDetailResponse, error) {
	members := make([]dto.TeamMemberInfo, 0, 3)

	// 队长
	var captain models.User
	if err := config.DB.Select("id, username").First(&captain, team.CaptainID).Error; err == nil {
		members = append(members, dto.TeamMemberInfo{
			UserID:   captain.ID,
			Username: captain.Username,
			Role:     "captain",
		})
	}

	// 成员1
	if team.Member1ID != nil {
		var member1 models.User
		if err := config.DB.Select("id, username").First(&member1, *team.Member1ID).Error; err == nil {
			members = append(members, dto.TeamMemberInfo{
				UserID:   member1.ID,
				Username: member1.Username,
				Role:     "member1",
			})
		}
	}

	// 成员2
	if team.Member2ID != nil {
		var member2 models.User
		if err := config.DB.Select("id, username").First(&member2, *team.Member2ID).Error; err == nil {
			members = append(members, dto.TeamMemberInfo{
				UserID:   member2.ID,
				Username: member2.Username,
				Role:     "member2",
			})
		}
	}

	return &dto.TeamDetailResponse{
		TeamResponse: dto.TeamResponse{
			ID:          team.ID,
			TeamName:    team.TeamName,
			CaptainID:   team.CaptainID,
			CaptainName: team.CaptainName,
			Member1ID:   team.Member1ID,
			Member1Name: team.Member1Name,
			Member2ID:   team.Member2ID,
			Member2Name: team.Member2Name,
			TeamTrack:   team.TeamTrack,
			SchoolID:    team.SchoolID,
			SchoolName:  team.SchoolName,
			TeamScore:   team.TeamScore,
			MemberCount: team.MemberCount,
			Status:      team.Status,
			CreatedAt:   team.CreatedAt,
			UpdatedAt:   team.UpdatedAt,
		},
		Members: members,
	}, nil
}

// UpdateTeam 更新团队信息（队长）
func (s *TeamService) UpdateTeam(userID int64, req *dto.UpdateTeamRequest) error {
	// 查找用户所在的团队
	var team models.Team
	if err := config.DB.Where("captain_id = ? AND status = 'active'", userID).First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("您不是任何团队的队长")
		}
		return err
	}

	updates := make(map[string]interface{})

	if req.TeamName != nil {
		// 检查新团队名称是否已被使用
		var count int64
		if err := config.DB.Model(&models.Team{}).
			Where("team_name = ? AND id != ?", *req.TeamName, team.ID).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("团队名称已存在")
		}
		updates["team_name"] = *req.TeamName
	}

	if req.TeamPassword != nil {
		// 加密新密码
		if err := team.SetPassword(*req.TeamPassword); err != nil {
			return err
		}
		updates["team_password"] = team.TeamPassword
	}

	if len(updates) > 0 {
		if err := config.DB.Model(&team).Updates(updates).Error; err != nil {
			return err
		}
	}

	return nil
}

// TransferCaptain 转让队长（队长）
func (s *TeamService) TransferCaptain(userID int64, req *dto.TransferCaptainRequest) error {
	// 查找用户所在的团队
	var team models.Team
	if err := config.DB.Where("captain_id = ? AND status = 'active'", userID).First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("您不是任何团队的队长")
		}
		return err
	}

	// 检查新队长是否是团队成员
	if !team.IsMember(req.NewCaptainID) {
		return errors.New("新队长必须是团队成员")
	}

	// 获取新队长信息
	var newCaptain models.User
	if err := config.DB.Select("id, username").First(&newCaptain, req.NewCaptainID).Error; err != nil {
		return errors.New("新队长不存在")
	}

	// 更新队长
	updates := map[string]interface{}{
		"captain_id":   newCaptain.ID,
		"captain_name": newCaptain.Username,
	}

	// 将原队长移到成员位置
	if team.Member1ID != nil && *team.Member1ID == req.NewCaptainID {
		updates["member1_id"] = userID
		updates["member1_name"] = team.CaptainName
	} else if team.Member2ID != nil && *team.Member2ID == req.NewCaptainID {
		updates["member2_id"] = userID
		updates["member2_name"] = team.CaptainName
	}

	if err := config.DB.Model(&team).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// RemoveMember 移除成员（队长）
func (s *TeamService) RemoveMember(userID int64, req *dto.RemoveMemberRequest) error {
	// 查找用户所在的团队
	var team models.Team
	if err := config.DB.Where("captain_id = ? AND status = 'active'", userID).First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("您不是任何团队的队长")
		}
		return err
	}

	// 不能移除队长自己
	if req.MemberID == userID {
		return errors.New("不能移除队长自己")
	}

	// 检查是否是团队成员
	if !team.IsMember(req.MemberID) {
		return errors.New("该用户不是团队成员")
	}

	// 移除成员
	updates := make(map[string]interface{})
	if team.Member1ID != nil && *team.Member1ID == req.MemberID {
		updates["member1_id"] = nil
		updates["member1_name"] = nil
	} else if team.Member2ID != nil && *team.Member2ID == req.MemberID {
		updates["member2_id"] = nil
		updates["member2_name"] = nil
	}
	updates["member_count"] = team.MemberCount - 1

	if err := config.DB.Model(&team).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// DisbandTeam 解散团队（队长）
func (s *TeamService) DisbandTeam(userID int64) error {
	// 查找用户所在的团队
	var team models.Team
	if err := config.DB.Where("captain_id = ? AND status = 'active'", userID).First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("您不是任何团队的队长")
		}
		return err
	}

	// 更新团队状态为已解散
	if err := config.DB.Model(&team).Update("status", "disbanded").Error; err != nil {
		return err
	}

	return nil
}

// GetTeamList 获取团队列表
func (s *TeamService) GetTeamList(req *dto.TeamListRequest) (*dto.TeamListResponse, error) {
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
	query := config.DB.Model(&models.Team{})

	// 搜索条件
	if req.Search != "" {
		query = query.Where("team_name LIKE ? OR captain_name LIKE ?",
			"%"+req.Search+"%", "%"+req.Search+"%")
	}

	// 筛选条件
	if req.TeamTrack != "" {
		query = query.Where("team_track = ?", req.TeamTrack)
	}
	if req.SchoolID != nil {
		query = query.Where("school_id = ?", *req.SchoolID)
	}
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
	var teams []models.Team
	if err := query.Offset(offset).Limit(req.Limit).Find(&teams).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]dto.TeamResponse, 0, len(teams))
	for _, team := range teams {
		list = append(list, dto.TeamResponse{
			ID:          team.ID,
			TeamName:    team.TeamName,
			CaptainID:   team.CaptainID,
			CaptainName: team.CaptainName,
			Member1ID:   team.Member1ID,
			Member1Name: team.Member1Name,
			Member2ID:   team.Member2ID,
			Member2Name: team.Member2Name,
			TeamTrack:   team.TeamTrack,
			SchoolID:    team.SchoolID,
			SchoolName:  team.SchoolName,
			TeamScore:   team.TeamScore,
			MemberCount: team.MemberCount,
			Status:      team.Status,
			CreatedAt:   team.CreatedAt,
			UpdatedAt:   team.UpdatedAt,
		})
	}

	return &dto.TeamListResponse{
		Total: int(total),
		Page:  req.Page,
		Limit: req.Limit,
		List:  list,
	}, nil
}

// UpdateTeamStatus 更新团队状态（管理员）
func (s *TeamService) UpdateTeamStatus(teamID int64, status string) error {
	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("团队不存在")
		}
		return err
	}

	if err := config.DB.Model(&team).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

// GetTeamRank 获取团队排名
func (s *TeamService) GetTeamRank(req *dto.TeamRankRequest) (*dto.TeamRankResponse, error) {
	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 50
	}

	// 构建查询
	query := config.DB.Model(&models.Team{}).Where("status = 'active'")

	// 筛选条件
	if req.TeamTrack != "" {
		query = query.Where("team_track = ?", req.TeamTrack)
	}
	if req.SchoolID != nil {
		query = query.Where("school_id = ?", *req.SchoolID)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 按分数排序，分数相同按最后解题时间排序
	query = query.Order("team_score DESC, updated_at ASC")

	// 分页查询
	offset := (req.Page - 1) * req.Limit
	var teams []models.Team
	if err := query.Offset(offset).Limit(req.Limit).Find(&teams).Error; err != nil {
		return nil, err
	}

	// 构建排名列表
	list := make([]dto.TeamRankItem, 0, len(teams))
	for i, team := range teams {
		// 查询解题数量
		var solveCount int64
		config.DB.Model(&models.Solve{}).Where("team_id = ?", team.ID).Count(&solveCount)

		// 查询最后解题时间
		var lastSolve models.Solve
		var lastSolveAt *time.Time
		if err := config.DB.Where("team_id = ?", team.ID).
			Order("solving_time DESC").
			First(&lastSolve).Error; err == nil {
			lastSolveAt = &lastSolve.SolvingTime
		}

		list = append(list, dto.TeamRankItem{
			Rank:        offset + i + 1,
			TeamID:      team.ID,
			TeamName:    team.TeamName,
			TeamScore:   team.TeamScore,
			TeamTrack:   team.TeamTrack,
			SchoolName:  team.SchoolName,
			MemberCount: team.MemberCount,
			SolveCount:  int(solveCount),
			LastSolveAt: lastSolveAt,
		})
	}

	return &dto.TeamRankResponse{
		Total: int(total),
		Page:  req.Page,
		Limit: req.Limit,
		List:  list,
	}, nil
}
