package postgres

import (
	"context"
	"database/sql"
	"time"

	"trading-ace/internal/domain"
)

type CampaignRepository struct {
	db *sql.DB
}

func NewCampaignRepository(db *sql.DB) *CampaignRepository {
	return &CampaignRepository{db: db}
}

func (r *CampaignRepository) GetActiveCampaign(ctx context.Context) (*domain.Campaign, error) {
	query := `
        SELECT id, start_time, end_time, status 
        FROM campaigns 
        WHERE status = 'ACTIVE' 
        AND start_time <= $1 
        AND end_time > $1 
        LIMIT 1
    `

	campaign := &domain.Campaign{}
	err := r.db.QueryRowContext(ctx, query, time.Now()).Scan(
		&campaign.ID,
		&campaign.StartTime,
		&campaign.EndTime,
		&campaign.Status,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return campaign, err
}
