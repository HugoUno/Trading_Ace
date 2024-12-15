package domain

import (
	"time"
)

type Campaign struct {
	ID        int64     `json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"` // 即將開始, 進行中, 已結束
}

type Task struct {
	ID          int64     `json:"id"`
	CampaignID  int64     `json:"campaign_id"`
	Type        string    `json:"type"` // 入門任務, 分享池任務
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	PoolAddress string    `json:"pool_address"`
	PointsPool  int64     `json:"points_pool"`
}

type UserTask struct {
	ID        int64     `json:"id"`
	UserID    string    `json:"user_id"` // 以太坊地址
	TaskID    int64     `json:"task_id"`
	Status    string    `json:"status"` // 進行中, 已完成
	Amount    float64   `json:"amount"`
	Points    int64     `json:"points"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRanking 用戶排名資訊
type UserRanking struct {
	UserID      string    `json:"user_id"`
	TotalPoints int64     `json:"total_points"`
	Rank        int       `json:"rank"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserVolume 用戶交易量資訊
type UserVolume struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}
