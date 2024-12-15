package postgres

import (
	"context"
	"database/sql"

	"trading-ace/internal/domain"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) GetTasksByCampaignID(ctx context.Context, campaignID int64) ([]*domain.Task, error) {
	query := `
        SELECT id, campaign_id, type, start_time, end_time, pool_address, points_pool 
        FROM tasks 
        WHERE campaign_id = $1
    `

	rows, err := r.db.QueryContext(ctx, query, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		task := &domain.Task{}
		err := rows.Scan(
			&task.ID,
			&task.CampaignID,
			&task.Type,
			&task.StartTime,
			&task.EndTime,
			&task.PoolAddress,
			&task.PointsPool,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (r *TaskRepository) GetTaskByID(ctx context.Context, taskID int64) (*domain.Task, error) {
	query := `
        SELECT id, campaign_id, type, start_time, end_time, pool_address, points_pool 
        FROM tasks 
        WHERE id = $1
    `

	task := &domain.Task{}
	err := r.db.QueryRowContext(ctx, query, taskID).Scan(
		&task.ID,
		&task.CampaignID,
		&task.Type,
		&task.StartTime,
		&task.EndTime,
		&task.PoolAddress,
		&task.PointsPool,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return task, err
}
