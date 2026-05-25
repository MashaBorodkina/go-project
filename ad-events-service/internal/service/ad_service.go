package service

import (
	"ad-events-service/internal/model"
	"ad-events-service/internal/repository"
	"context"
	"fmt"
)

type AdService struct {
	BanRepo  *repository.BannerRepository
	CampRepo *repository.Repository
}

func NewAdService(banRepo *repository.BannerRepository, campRepo *repository.Repository) *AdService {
	return &AdService{
		BanRepo:  banRepo,
		CampRepo: campRepo,
	}
}

func (s *AdService) GetBannerForDisplay(ctx context.Context) (*model.Banner, error) {
	camps, err := s.CampRepo.GetAllCampaigns(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all banners: %w", err)
	}

	campIdFound := ""

	for _, camp := range camps {
		if camp.Budget > 0 && camp.Status == "active" {
			campIdFound = camp.ID.String()
			break
		}
	}

	if campIdFound == "" {
		return nil, fmt.Errorf("no suitable campaign found")
	}

	bans, err := s.BanRepo.GetBannersByCampaignId(ctx, campIdFound)

	if err != nil {
		return nil, fmt.Errorf("failed to get banners by campaign ID: %w", err)
	}

	for _, ban := range bans {
		if ban.IsActive {
			return ban, nil
		}
	}
	return nil, fmt.Errorf("no active banner found for display")
}
