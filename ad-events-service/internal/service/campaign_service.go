package service

import (
	"ad-events-service/internal/dto"
	"ad-events-service/internal/model"
	"ad-events-service/internal/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CampaignService struct {
	campRepo *repository.Repository
}

func NewCampaignService(campRepo *repository.Repository) *CampaignService {
	return &CampaignService{campRepo: campRepo}
}

func (s *CampaignService) GetCampaignByID(ctx context.Context, campID string) (*model.Campaign, error) {
	if campID == "" {
		return nil, fmt.Errorf("campaign ID cannot be empty")
	}
	
	camp, err := s.campRepo.GetCampaignByID(ctx, campID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign by ID: %w", err)
	}

	if camp == nil {
		return nil, fmt.Errorf("no campaign found for campaign ID: %s", campID)
	}
	return camp, nil
}

func (s *CampaignService) GetAllCampaigns(ctx context.Context) ([]*model.Campaign, error) {
	camps, err := s.campRepo.GetAllCampaigns(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all campaigns: %w", err)
	}	
	return camps, nil
}

func ValidateCampaign(camp *model.Campaign) error {
	switch {
	case camp.Name == "":
		return fmt.Errorf("campaign name cannot be empty")
	case len(camp.Name) > 100 || len(camp.Name) < 1:
		return fmt.Errorf("campaign name must be between 1 and 100 characters")
	case camp.Budget <= 0:
		return fmt.Errorf("campaign budget must be a positive number")
	case isValidStatus(camp.Status):
		return fmt.Errorf("invalid campaign status: %s", camp.Status)
	case camp.ID == uuid.Nil:
		return fmt.Errorf("campaign ID cannot be empty")
	}
	return nil
}

func isValidStatus(status string) bool {
	switch status {
	case "active", "paused", "archived":
		return true
	default:
		return false
	}
}

func (s *CampaignService) CreateCampaign(ctx context.Context, camp *model.Campaign) error {
	if err := ValidateCampaign(camp); err != nil {
		return fmt.Errorf("campaign validation failed: %w", err)
	}
	if err := s.campRepo.CreateCampaign(ctx, camp); err != nil {
		return fmt.Errorf("failed to create campaign: %w", err)
	}
	return nil
}

func (s *CampaignService) UpdateCampaign(ctx context.Context, camp *model.Campaign) error {
	if err := ValidateCampaign(camp); err != nil {
		return fmt.Errorf("campaign validation failed: %w", err)
	}
	if err := s.campRepo.UpdateCampaign(ctx, camp); err != nil {
		return fmt.Errorf("failed to update campaign: %w", err)
	}
	return nil
}

func (s *CampaignService) DeleteCampaign(ctx context.Context, campID string) error {
	if campID == "" {
		return fmt.Errorf("campaign ID cannot be empty")
	}
	if err := s.campRepo.DeleteCampaign(ctx, campID); err != nil {
		return fmt.Errorf("failed to delete campaign: %w", err)
	}
	return nil
}

// PatchCampaign allows partial updates to a campaign. Only non-nil fields in the request will be updated.
func (s *CampaignService) PatchCampaign(ctx context.Context, campID string, req *dto.PatchCampaignRequest) (*model.Campaign, error) {
	switch {
	case campID == "":
		return nil, fmt.Errorf("campaign ID cannot be empty")
	case req == nil:
		return nil, fmt.Errorf("patch request cannot be nil")
	case req.Name != nil && (len(*req.Name) > 100 || len(*req.Name) < 1):
		return nil, fmt.Errorf("campaign name must be between 1 and 100 characters")
	case req.Budget != nil && *req.Budget <= 0:
		return nil, fmt.Errorf("campaign budget must be a positive number")
	case req.Status != nil && !isValidStatus(*req.Status):
		return nil, fmt.Errorf("invalid campaign status: %s", *req.Status)
	}

	campaign, err := s.campRepo.GetCampaignByID(ctx, campID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign by ID: %w", err)
	}
	if req.Name != nil {
		campaign.Name = *req.Name
	}
	if req.Budget != nil {
		campaign.Budget = *req.Budget
	}
	if req.Status != nil {
		campaign.Status = *req.Status
	}
	if err := s.campRepo.UpdateCampaign(ctx, campaign); err != nil {
		return nil, fmt.Errorf("failed to update campaign: %w", err)
	}
	return campaign, nil
}
