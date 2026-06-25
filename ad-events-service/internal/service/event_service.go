package service

import (
	"context"
	"fmt"

	"ad-events-service/internal/apperrors"
	"ad-events-service/internal/model"
	"ad-events-service/internal/repository"
)

type EventService struct {
	EventRepo *repository.EventRepository
	BanRepo   *repository.BannerRepository
	CampRepo  *repository.Repository
}

func NewEventService(
	eventRepo *repository.EventRepository,
	campRepo *repository.Repository,
	banRepo *repository.BannerRepository,
) *EventService {
	return &EventService{
		EventRepo: eventRepo,
		CampRepo:  campRepo,
		BanRepo:   banRepo,
	}
}

func (s *EventService) TrackEvent(
	ctx context.Context,
	ban_id string,
	eventType string,
	ip string,
	userAgent string,
) error {
	if ban_id == "" {
		return fmt.Errorf("banner ID cannot be empty")
	}
	ban, err := s.BanRepo.GetBannerByID(ctx, ban_id)
	if err != nil {
		return fmt.Errorf("failed to get banner by ID: %w", err)
	}

	if ban == nil {
		return apperrors.ErrBannerNotFound
	}
	if !ban.IsActive {
		return apperrors.ErrBannerInactive
	}

	camp, err := s.CampRepo.GetCampaignByID(ctx, ban.CampaignID.String())
	if err != nil {
		return fmt.Errorf("failed to get campaign by banner ID: %w", err)
	}

	if camp == nil {
		return apperrors.ErrCampaignNotFound
	}

	if camp.Status != model.CampaignStatusActive {
		return apperrors.ErrCampaignInactive
	}

	if eventType != model.EventTypeImpression && eventType != model.EventTypeClick {
		return apperrors.ErrInvalidEventType
	}

	event := &model.Event{
		Banner_ID:  ban.ID,
		Type:       eventType,
		Ip:         ip,
		User_Agent: userAgent,
	}
	if err := s.EventRepo.CreateEvent(ctx, event); err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	return nil
}

func (s *EventService) TrackImpression(ctx context.Context, ban_id string, ip string, userAgent string) error {
	return s.TrackEvent(ctx, ban_id, model.EventTypeImpression, ip, userAgent)
}

func (s *EventService) TrackClick(ctx context.Context, ban_id string, ip string, userAgent string) error {
	return s.TrackEvent(ctx, ban_id, model.EventTypeClick, ip, userAgent)
}

func (s *EventService) GetEventByID(ctx context.Context, ID string) (*model.Event, error) {
	if ID == "" {
		return nil, fmt.Errorf("event ID cannot be empty")
	}
	event, err := s.EventRepo.GetEventByID(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event by ID: %w", err)
	}
	if event == nil {
		return nil, apperrors.ErrEventNotFound
	}

	return event, nil
}

func (s *EventService) GetAllEvents(ctx context.Context) ([]*model.Event, error) {
	events, err := s.EventRepo.GetAllEvents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all events: %w", err)
	}
	return events, nil
}

func (s *EventService) GetEventsByBannerID(
	ctx context.Context,
	bannerID string,
	eventType string,
	limit int,
	offset int,
) ([]*model.Event, error) {
	if bannerID == "" {
		return nil, fmt.Errorf("banner ID cannot be empty")
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	if eventType != "" && eventType != model.EventTypeImpression && eventType != model.EventTypeClick {
		return nil, apperrors.ErrInvalidEventType
	}
	events, err := s.EventRepo.GetEventsByBannerID(ctx, bannerID, eventType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by banner ID: %w", err)
	}
	return events, nil
}

func (s *EventService) GetEventsByCampaignID(
	ctx context.Context,
	campaignID string,
	eventType string,
	limit int,
	offset int,
) ([]*model.Event, error) {
	if campaignID == "" {
		return nil, fmt.Errorf("campaign ID cannot be empty")
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	if eventType != "" && eventType != model.EventTypeImpression && eventType != model.EventTypeClick {
		return nil, apperrors.ErrInvalidEventType
	}
	events, err := s.EventRepo.GetEventsByCampaignID(ctx, campaignID, eventType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by campaign ID: %w", err)
	}
	return events, nil
}
