package service

import (
	"ad-events-service/internal/apperrors"
	"ad-events-service/internal/model"
	"ad-events-service/internal/repository"
	"context"
	"time"
)

type StatsService struct {
	StatsRepo *repository.StatsRepository
}

func NewStatsService(statsRepo *repository.StatsRepository) *StatsService {
	return &StatsService{
		StatsRepo: statsRepo,
	}
}

func (s *StatsService) GetBannerStatsByID(ctx context.Context, bannerID string, from time.Time, to time.Time) (*model.BannerStats, error) {
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
	return stats, nil
}
