package service

import (
	"context"
	"fmt"
	"time"

	"trading-ace/internal/domain"
	"trading-ace/pkg/ethereum"
)

type CampaignService struct {
	campaignRepo CampaignRepository
	taskRepo     TaskRepository
	userTaskRepo UserTaskRepository
	ethClient    *ethereum.Client
}

type CampaignRepository interface {
	GetActiveCampaign(ctx context.Context) (*domain.Campaign, error)
}

type TaskRepository interface {
	GetTasksByCampaignID(ctx context.Context, campaignID int64) ([]*domain.Task, error)
	GetTaskByID(ctx context.Context, taskID int64) (*domain.Task, error)
}

type UserTaskRepository interface {
	CreateUserTask(ctx context.Context, task *domain.UserTask) error
	UpdateUserTask(ctx context.Context, task *domain.UserTask) error
	GetUserTasks(ctx context.Context, userID string) ([]*domain.UserTask, error)
	GetUserVolumesByTaskID(ctx context.Context, taskID int64) ([]*domain.UserVolume, error)
	UpdateUserTaskPoints(ctx context.Context, taskID int64, userID string, points int64) error
	GetUserRankings(ctx context.Context, tasks []*domain.Task) ([]*domain.UserRanking, error)
}

// 分享池積分計���結構
type SharePoolResult struct {
	UserID string
	Volume float64
	Share  float64
	Points int64
}

func NewCampaignService(
	campaignRepo CampaignRepository,
	taskRepo TaskRepository,
	userTaskRepo UserTaskRepository,
	ethClient *ethereum.Client,
) *CampaignService {
	return &CampaignService{
		campaignRepo: campaignRepo,
		taskRepo:     taskRepo,
		userTaskRepo: userTaskRepo,
		ethClient:    ethClient,
	}
}

// 處理 Swap 事件
func (s *CampaignService) HandleSwapEvent(ctx context.Context, event *ethereum.SwapEvent) error {
	campaign, err := s.campaignRepo.GetActiveCampaign(ctx)
	if err != nil {
		return err
	}

	tasks, err := s.taskRepo.GetTasksByCampaignID(ctx, campaign.ID)
	if err != nil {
		return err
	}

	// 計算交易金額
	amount := calculateTradeAmount(event)

	// 更新用戶任務
	for _, task := range tasks {
		userTask := &domain.UserTask{
			UserID:    event.Sender.Hex(),
			TaskID:    task.ID,
			Amount:    amount,
			UpdatedAt: time.Now(),
		}

		// 入門任務檢查
		if task.Type == "ONBOARDING" && amount >= 1000 {
			userTask.Status = "COMPLETED"
			userTask.Points = 100
		}

		if err := s.userTaskRepo.CreateUserTask(ctx, userTask); err != nil {
			return err
		}
	}

	return nil
}

// 計算交易金額（這裡需要根據實際代幣格來計算）
func calculateTradeAmount(event *ethereum.SwapEvent) float64 {
	// 計算 USDC 的交易量（假設 USDC 是 token1）
	if event.Amount1In.Sign() > 0 {
		return float64(event.Amount1In.Uint64()) / 1e6 // USDC 有 6 位小數
	}
	if event.Amount1Out.Sign() > 0 {
		return float64(event.Amount1Out.Uint64()) / 1e6
	}
	return 0
}

// 獲取用戶任務狀態
func (s *CampaignService) GetUserTaskStatus(ctx context.Context, userID string) ([]*domain.UserTask, error) {
	return s.userTaskRepo.GetUserTasks(ctx, userID)
}

// 計算分享池積分
func (s *CampaignService) CalculateSharePoolPoints(ctx context.Context, taskID int64) error {
	// 1. 獲取任務資訊
	task, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("獲取任務失敗: %v", err)
	}

	// 2. 獲取該任務期間內的所有用戶交易量
	userVolumes, err := s.userTaskRepo.GetUserVolumesByTaskID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("獲取用戶交易量失敗: %v", err)
	}

	// 3. 計算總交易量
	var totalVolume float64
	for _, volume := range userVolumes {
		totalVolume += volume.Amount
	}

	// 4. 計算每個用戶的分享比例和積分
	results := make([]*SharePoolResult, 0, len(userVolumes))
	for _, volume := range userVolumes {
		// 檢查用戶是否完成入門任務
		completed, err := s.hasCompletedOnboarding(ctx, volume.UserID)
		if err != nil {
			return fmt.Errorf("檢查入門任務失敗: %v", err)
		}
		if !completed {
			continue
		}

		// 計算分享比例
		share := volume.Amount / totalVolume
		points := int64(share * float64(task.PointsPool))

		results = append(results, &SharePoolResult{
			UserID: volume.UserID,
			Volume: volume.Amount,
			Share:  share,
			Points: points,
		})
	}

	// 5. 更新用戶積分
	for _, result := range results {
		err := s.userTaskRepo.UpdateUserTaskPoints(ctx, taskID, result.UserID, result.Points)
		if err != nil {
			return fmt.Errorf("更新用戶積分失敗: %v", err)
		}
	}

	return nil
}

// 檢查用戶是否完成入門任務
func (s *CampaignService) hasCompletedOnboarding(ctx context.Context, userID string) (bool, error) {
	tasks, err := s.userTaskRepo.GetUserTasks(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, task := range tasks {
		if task.Status == "COMPLETED" {
			return true, nil
		}
	}
	return false, nil
}

// GetLeaderboard 獲取排行榜
func (s *CampaignService) GetLeaderboard(ctx context.Context) ([]*domain.UserRanking, error) {
	campaign, err := s.campaignRepo.GetActiveCampaign(ctx)
	if err != nil {
		return nil, fmt.Errorf("獲取活動失敗: %v", err)
	}

	tasks, err := s.taskRepo.GetTasksByCampaignID(ctx, campaign.ID)
	if err != nil {
		return nil, fmt.Errorf("獲取任務失敗: %v", err)
	}

	// 獲取所有已完成任務的用戶積分
	rankings, err := s.userTaskRepo.GetUserRankings(ctx, tasks)
	if err != nil {
		return nil, fmt.Errorf("獲取排名失敗: %v", err)
	}

	return rankings, nil
}
