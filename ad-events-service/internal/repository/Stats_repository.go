package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"ad-events-service/internal/model"
)

type StatsRepository struct {
	db *pgxpool.Pool
}

func NewStatsRepository(db *pgxpool.Pool) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) GetBannerStatsByID(
	ctx context.Context,
	bannerID string,
	from time.Time,
	to time.Time,
) (*model.BannerStats, error) {
	var stats model.BannerStats
	query := `SELECT type, COUNT(*) as count 
	FROM events 
	WHERE banner_id = $1
	`
	args := []any{bannerID}

	if !from.IsZero() && !to.IsZero() {
		query += ` AND created_at >= $2 AND created_at <= $3`
		args = append(args, from, to)
	}

	query += ` group by type`

	rows, err := r.db.Query(ctx, query, args...)
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
		case model.EventTypeImpression:
			stats.Impressions = count
		case model.EventTypeClick:
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

func (r *StatsRepository) GetCampaignStatsByID(
	ctx context.Context,
	campaignID string,
	from time.Time,
	to time.Time,
) (*model.CampaignStats, error) {
	var stats model.CampaignStats
	query := `SELECT c.name, e.type, COUNT(*) as count, c.budget 
	FROM events e 
	JOIN banners b ON e.banner_id = b.id 
	JOIN campaigns c ON b.campaign_id = c.id 
	WHERE b.campaign_id = $1 
	`
	args := []any{campaignID}
	if !from.IsZero() && !to.IsZero() {
		query += ` AND e.created_at >= $2 AND e.created_at <= $3`
		args = append(args, from, to)
	}

	query += ` group by c.name, e.type, c.budget`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign stats by ID: %w", err)
	}
	defer rows.Close()

	id, err := uuid.Parse(campaignID)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign ID format: %w", err)
	}

	stats.CampaignID = id

	found := false

	for rows.Next() {
		found = true
		var eventType string
		var count int
		if err := rows.Scan(&stats.CampaignName, &eventType, &count, &stats.Budget); err != nil {
			return nil, fmt.Errorf("failed to scan campaign stats: %w", err)
		}
		switch eventType {
		case model.EventTypeImpression:
			stats.Impressions = count

		case model.EventTypeClick:
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

func (r *StatsRepository) GetDailyStats(
	ctx context.Context,
	campaignID string,
	from time.Time,
	to time.Time,
) ([]*model.DailyStats, error) {
	query := `Select Date(e.created_at) as date, e.type, Count(*) as count
	 From events e
	 Join banners b On e.banner_id = b.id
	 Where b.campaign_id = $1 
	 `
	args := []any{campaignID}

	if !from.IsZero() && !to.IsZero() {
		query += ` AND e.created_at >= $2 AND e.created_at <= $3`
		args = append(args, from, to)
	}

	query += ` Group By Date(e.created_at), e.type`

	rows, err := r.db.Query(ctx, query, args...)
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
		case model.EventTypeImpression:
			statsMap[dateKey].Impressions = count
		case model.EventTypeClick:
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
