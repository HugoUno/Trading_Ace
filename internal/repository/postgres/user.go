package postgres

import (
	"context"
	"database/sql"

	"trading-ace/internal/domain"

	"github.com/lib/pq"
)

type UserTaskRepository struct {
	db *sql.DB
}

func NewUserTaskRepository(db *sql.DB) *UserTaskRepository {
	return &UserTaskRepository{db: db}
}

func (r *UserTaskRepository) CreateUserTask(ctx context.Context, task *domain.UserTask) error {
	query := `
        INSERT INTO user_tasks (user_id, task_id, status, amount, points, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (user_id, task_id) DO UPDATE
        SET amount = user_tasks.amount + EXCLUDED.amount,
            points = CASE 
                WHEN user_tasks.status = 'COMPLETED' THEN user_tasks.points 
                ELSE EXCLUDED.points 
            END,
            status = CASE 
                WHEN user_tasks.status = 'COMPLETED' THEN user_tasks.status 
                ELSE EXCLUDED.status 
            END,
            updated_at = EXCLUDED.updated_at
    `

	_, err := r.db.ExecContext(ctx, query,
		task.UserID,
		task.TaskID,
		task.Status,
		task.Amount,
		task.Points,
		task.UpdatedAt,
	)

	return err
}

func (r *UserTaskRepository) UpdateUserTask(ctx context.Context, task *domain.UserTask) error {
	query := `
        UPDATE user_tasks 
        SET status = $1, amount = $2, points = $3, updated_at = $4
        WHERE id = $5
    `

	_, err := r.db.ExecContext(ctx, query,
		task.Status,
		task.Amount,
		task.Points,
		task.UpdatedAt,
		task.ID,
	)

	return err
}

func (r *UserTaskRepository) GetUserTasks(ctx context.Context, userID string) ([]*domain.UserTask, error) {
	query := `
        SELECT id, user_id, task_id, status, amount, points, updated_at
        FROM user_tasks
        WHERE user_id = $1
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.UserTask
	for rows.Next() {
		task := &domain.UserTask{}
		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.TaskID,
			&task.Status,
			&task.Amount,
			&task.Points,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (r *UserTaskRepository) GetUserVolumesByTaskID(ctx context.Context, taskID int64) ([]*domain.UserVolume, error) {
	query := `
        SELECT user_id, SUM(amount) as total_amount
        FROM user_tasks
        WHERE task_id = $1
        GROUP BY user_id
    `

	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var volumes []*domain.UserVolume
	for rows.Next() {
		volume := &domain.UserVolume{}
		err := rows.Scan(&volume.UserID, &volume.Amount)
		if err != nil {
			return nil, err
		}
		volumes = append(volumes, volume)
	}

	return volumes, rows.Err()
}

func (r *UserTaskRepository) UpdateUserTaskPoints(ctx context.Context, taskID int64, userID string, points int64) error {
	query := `
        UPDATE user_tasks 
        SET points = $1, status = 'COMPLETED'
        WHERE task_id = $2 AND user_id = $3
    `

	_, err := r.db.ExecContext(ctx, query, points, taskID, userID)
	return err
}

func (r *UserTaskRepository) GetUserRankings(ctx context.Context, tasks []*domain.Task) ([]*domain.UserRanking, error) {
	// 構建任務 ID 列表
	taskIDs := make([]int64, len(tasks))
	for i, task := range tasks {
		taskIDs[i] = task.ID
	}

	query := `
        WITH RankedUsers AS (
            SELECT 
                user_id,
                SUM(points) as total_points,
                MAX(updated_at) as last_updated,
                RANK() OVER (ORDER BY SUM(points) DESC) as rank
            FROM user_tasks
            WHERE task_id = ANY($1)
                AND status = 'COMPLETED'
            GROUP BY user_id
        )
        SELECT user_id, total_points, rank, last_updated
        FROM RankedUsers
        ORDER BY rank ASC
    `

	rows, err := r.db.QueryContext(ctx, query, pq.Array(taskIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rankings []*domain.UserRanking
	for rows.Next() {
		ranking := &domain.UserRanking{}
		err := rows.Scan(
			&ranking.UserID,
			&ranking.TotalPoints,
			&ranking.Rank,
			&ranking.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rankings = append(rankings, ranking)
	}

	return rankings, rows.Err()
}
