package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"ad-events-service/internal/apperrors"
	"ad-events-service/internal/dto"
	"ad-events-service/internal/model"
	"ad-events-service/internal/repository"
)

type BannerService struct {
	BanRepo  *repository.BannerRepository
	CampRepo *repository.Repository
}

func NewBannerService(BanRepo *repository.BannerRepository, CampRepo *repository.Repository) *BannerService {
	return &BannerService{BanRepo: BanRepo, CampRepo: CampRepo}
}

func (s *BannerService) GetBannerByID(ctx context.Context, banID string) (*model.Banner, error) {
	if banID == "" {
		return nil, apperrors.ErrInvalidBannerID
	}
	ban, err := s.BanRepo.GetBannerByID(ctx, banID)
	if err != nil {
		return nil, err
	}

	return ban, nil
}

func (s *BannerService) GetAllBannersByCampaignId(ctx context.Context, campaignId string) ([]*model.Banner, error) {
	if campaignId == "" {
		return nil, fmt.Errorf("campaign ID cannot be empty")
	}
	if _, err := uuid.Parse(campaignId); err != nil {
		return nil, fmt.Errorf("invalid campaign ID format: %w", err)
	}
	if _, err := s.CampRepo.GetCampaignByID(ctx, campaignId); err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}
	bans, err := s.BanRepo.GetAllBannersByCampaignId(ctx, campaignId)
	if err != nil {
		return nil, fmt.Errorf("failed to get all banners: %w", err)
	}
	if len(bans) == 0 {
		return nil, apperrors.ErrBannerNotFound
	}

	return bans, nil
}

func ValidateBanner(ban *model.Banner) error {
	switch {
	case len(ban.Title) > 200 || len(ban.Title) < 1:
		return fmt.Errorf("banner title must be between 1 and 200 characters")
	case ban.ImageUrl == "":
		return fmt.Errorf("banner image URL cannot be empty")
	case len(ban.ImageUrl) > 500 || len(ban.ImageUrl) < 1:
		return fmt.Errorf("banner image URL must be between 1 and 500 characters")
	case ban.CampaignID == uuid.Nil:
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

func (s *BannerService) UpdateBanner(ctx context.Context, banner *model.Banner) error {
	if err := ValidateBanner(banner); err != nil {
		return fmt.Errorf("banner validation failed: %w", err)
	}
	if err := s.BanRepo.UpdateBanner(ctx, banner); err != nil {
		return fmt.Errorf("failed to update banner: %w", err)
	}

	return nil
}

func (s *BannerService) PatchBanner(
	ctx context.Context,
	banID string,
	req *dto.BannerPatchRequest,
) (*model.Banner, error) {
	switch {
	case banID == "":
		return nil, fmt.Errorf("banner ID cannot be empty")
	case req == nil:
		return nil, fmt.Errorf("patch request cannot be nil")
	case req.Title != nil && (len(*req.Title) > 200 || len(*req.Title) < 1):
		return nil, fmt.Errorf("banner title must be between 1 and 200 characters")
	case req.ImageUrl != nil && (len(*req.ImageUrl) > 500 || len(*req.ImageUrl) < 1):
		return nil, fmt.Errorf("banner image URL must be between 1 and 500 characters")
	}

	banner, err := s.BanRepo.GetBannerByID(ctx, banID)
	if err != nil {
		return nil, fmt.Errorf("failed to get banner by ID: %w", err)
	}
	if req.Title != nil {
		banner.Title = *req.Title
	}
	if req.ImageUrl != nil {
		banner.ImageUrl = *req.ImageUrl
	}
	if req.IsActive != nil {
		banner.IsActive = *req.IsActive
	}
	if err := s.BanRepo.UpdateBanner(ctx, banner); err != nil {
		return nil, fmt.Errorf("failed to update banner: %w", err)
	}

	return banner, nil
}

func (s *BannerService) DeleteBanner(ctx context.Context, banID string) error {
	if banID == "" {
		return fmt.Errorf("banner ID cannot be empty")
	}
	if err := s.BanRepo.DeleteBanner(ctx, banID); err != nil {
		return fmt.Errorf("failed to delete banner: %w", err)
	}

	return nil
}
