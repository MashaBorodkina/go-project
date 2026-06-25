package service

import (
	"context"
	"math"
	"time"

	"ad-events-service/internal/apperrors"
	"ad-events-service/internal/model"
	"ad-events-service/internal/repository"
)

type StatsService struct {
	StatsRepo *repository.StatsRepository
	CampRepo  *repository.Repository
}

func NewStatsService(statsRepo *repository.StatsRepository, campRepo *repository.Repository) *StatsService {
	return &StatsService{
		StatsRepo: statsRepo,
		CampRepo:  campRepo,
	}
}

func CalculateCTR(clicks int, impressions int) float64 {
	if impressions == 0 {
		return 0.0
	}
	return math.Round((float64(clicks)/float64(impressions))*100*100) / 100
}

func (s *StatsService) GetBannerStatsByID(
	ctx context.Context,
	bannerID string,
	from time.Time,
	to time.Time,
) (*model.BannerStats, error) {
	if bannerID == "" {
		return nil, apperrors.ErrInvalidBannerID
	}
	if !from.IsZero() && !to.IsZero() && from.After(to) {
		return nil, apperrors.ErrInvalidDateRange
	}

	if from.IsZero() != to.IsZero() {
		return nil, apperrors.ErrBothDatesRequired
	}

	stats, err := s.StatsRepo.GetBannerStatsByID(ctx, bannerID, from, to)
	if err != nil {
		return nil, err
	}
	if stats == nil {
		return nil, apperrors.ErrEventNotFound
	}

	ctr := CalculateCTR(stats.Clicks, stats.Impressions)
	stats.CTR = ctr

	return stats, nil
}

func (s *StatsService) GetCampaignStatsByID(
	ctx context.Context,
	campaignID string,
	from time.Time,
	to time.Time,
) (*model.CampaignStats, error) {
	if campaignID == "" {
		return nil, apperrors.ErrInvalidCampaignID
	}
	if !from.IsZero() && !to.IsZero() && from.After(to) {
		return nil, apperrors.ErrInvalidDateRange
	}
	if from.IsZero() != to.IsZero() {
		return nil, apperrors.ErrBothDatesRequired
	}

	_, err := s.CampRepo.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	stats, err := s.StatsRepo.GetCampaignStatsByID(ctx, campaignID, from, to)
	if err != nil {
		return nil, err
	}

	ctr := CalculateCTR(stats.Clicks, stats.Impressions)
	stats.CTR = ctr

	return stats, nil
}

func (s *StatsService) GetDailyStats(
	ctx context.Context,
	campaignID string,
	from time.Time,
	to time.Time,
) ([]*model.DailyStats, error) {
	if campaignID == "" {
		return nil, apperrors.ErrInvalidCampaignID
	}
	if !from.IsZero() && !to.IsZero() && from.After(to) {
		return nil, apperrors.ErrInvalidDateRange
	}
	if from.IsZero() != to.IsZero() {
		return nil, apperrors.ErrBothDatesRequired
	}

	_, err := s.CampRepo.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	dailyStats, err := s.StatsRepo.GetDailyStats(ctx, campaignID, from, to)
	if err != nil {
		return nil, err
	}

	for _, stats := range dailyStats {
		stats.CTR = CalculateCTR(stats.Clicks, stats.Impressions)
	}

	return dailyStats, nil
}
