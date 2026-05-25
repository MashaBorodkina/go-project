package repository

import (
	"ad-events-service/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsRepository struct {
	db *pgxpool.Pool
}

func NewStatsRepository(db *pgxpool.Pool) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) GetBannerStatsByID(ctx context.Context, bannerID string) (*model.BannerStats, error) {
	var stats model.BannerStats
	query := "SELECT type, COUNT(*) as count FROM events WHERE banner_id = $1 group by type"
	rows, err := r.db.Query(ctx, query, bannerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get banner stats by ID: %w", err)
	}
	defer rows.Close()

	id, err := uuid.Parse(bannerID)
	if err != nil {
		return nil, fmt.Errorf("invalid banner ID format: %w", err)
	}

	stats.BannerID = id

	found := false

	for rows.Next() {
		found = true
		var eventType string
		var count int
		if err := rows.Scan(&eventType, &count); err != nil {
			return nil, fmt.Errorf("failed to scan banner stats: %w", err)
		}
		switch eventType {
		case "impression":
			stats.Impressions = count
		case "click":
			stats.Clicks = count
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over rows: %w", err)
	}
	if !found {
		stats.Impressions = 0
		stats.Clicks = 0
	}
	return &stats, nil
}

func (r *StatsRepository) GetCampaignStatsByID(ctx context.Context, campaignID string) (*model.CampaignStats, error) {
	var stats model.CampaignStats
	query := `SELECT b.campaign_id, c.name, e.type, COUNT(*) as count, c.budget 
	FROM events e 
	JOIN banners b ON e.banner_id = b.id 
	JOIN campaigns c ON b.campaign_id = c.id 
	WHERE b.campaign_id = $1 
	group by b.campaign_id, c.name, e.type, c.budget`
	rows, err := r.db.Query(ctx, query, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign stats by ID: %w", err)
	}
	defer rows.Close()

	found := false

	for rows.Next() {
		found = true
		var eventType string
		var count int
		if err := rows.Scan(&stats.CampaignID, &stats.CampaignName, &eventType, &count, &stats.Budget); err != nil {
			return nil, fmt.Errorf("failed to scan campaign stats: %w", err)
		}
		switch eventType {
		case "impression":
			stats.Impressions = count

		case "click":
			stats.Clicks = count
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over rows: %w", err)
	}

	if !found {
		stats.Impressions = 0
		stats.Clicks = 0
	}
	return &stats, nil
}

func (r *StatsRepository) GetDailyStats(ctx context.Context, campaignID string) ([]*model.DailyStats, error) {
	query := `Select b.campaign_id, e.type, Dte(e.created_at) as date, Count(*) as count
	 From events e
	 Join banners b On e.banner_id = b.id
	 Where b.campaign_id = $1 
	 Group By Date(e.created_at), e.type`

	rows, err := r.db.Query(ctx, query, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily stats: %w", err)
	}
	defer rows.Close()

	statsMap := make(map[string]*model.DailyStats)

	for rows.Next() {
		var date time.Time
		var eventType string
		var count int
		if err := rows.Scan(&date, &eventType, &count); err != nil {
			return nil, fmt.Errorf("failed to scan daily stats: %w", err)
		}

		dateKey := date.Format("2006-01-02")

		if statsMap[dateKey] == nil {
			statsMap[dateKey] = &model.DailyStats{Date: date}
		}
		switch eventType {
		case "impression":
			statsMap[dateKey].Impressions = count
		case "click":
			statsMap[dateKey].Clicks = count
		}

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over rows: %w", err)
	}
	var stats []*model.DailyStats
	for _, stat := range statsMap {
		stats = append(stats, stat)
	}
	return stats, nil
}
