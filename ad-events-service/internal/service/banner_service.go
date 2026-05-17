package service

import (
	"ad-events-service/internal/model"
	"ad-events-service/internal/repository"
	"context"
	"fmt"
)

type BannerService struct {
	BanRepo *repository.BannerRepository
}

func NewBannerService(BanRepo *repository.BannerRepository) *BannerService {
	return &BannerService{BanRepo: BanRepo}
}

func (s *BannerService) GetBannerByID(ctx context.Context, banID string) (*model.Banner, error) {
	if banID == "" {
		return nil, fmt.Errorf("banner ID cannot be empty")
	}
	ban, err := s.BanRepo.GetBannerByID(ctx, banID)
	if err != nil {
		return nil, fmt.Errorf("failed to get banner by ID: %w", err)
	}
	if ban == nil {
		return nil, fmt.Errorf("no banner found for banner ID: %s", banID)
	}
	return ban, nil
}

func (s *BannerService) GetAllBanners(ctx context.Context) ([]*model.Banner, error) {
	bans, err := s.BanRepo.GetAllBanners(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all banners: %w", err)
	}
	return bans, nil
}

func ValidateBanner(ban *model.Banner) error {
	switch {
	case ban.Title == "":
		return fmt.Errorf("banner title cannot be empty")
	case len(ban.Title) > 100 || len(ban.Title) < 1:
		return fmt.Errorf("banner title must be between 1 and 100 characters")
	case ban.ImageUrl == "":
		return fmt.Errorf("banner image URL cannot be empty")
	case len(ban.ImageUrl) > 255:
		return fmt.Errorf("banner image URL must be less than 255 characters")
	case ban.CampaignID == "":
		return fmt.Errorf("banner campaign ID cannot be empty")
	}
	return nil
}

func (s *BannerService) CreateBanner(ctx context.Context, ban *model.Banner) error {
	if err := ValidateBanner(ban); err != nil {
		return fmt.Errorf("invalid banner data: %w", err)
	}
	if err := s.BanRepo.CreateBanner(ctx, ban); err != nil {
		return fmt.Errorf("failed to create banner: %w", err)
	}
	return nil
}

func (s *BannerService) UpdateBanner(ctx context.Context, banner *model.Banner) error{
	if err := ValidateBanner(banner); err != nil {
		return fmt.Errorf("banner validation failed: %w", err)
	}
	if err := s.BanRepo.UpdateBanner(ctx, banner); err != nil {
		return fmt.Errorf("failed to update banner: %w", err)
	}	
	return nil
}