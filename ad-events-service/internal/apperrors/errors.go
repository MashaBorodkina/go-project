package apperrors

import "errors"

var (
	ErrBannerNotFound            = errors.New("banner or banners not found")
	ErrCampaignNotFound          = errors.New("campaign not found")
	ErrBannerInactive            = errors.New("banner is inactive")
	ErrCampaignInactive          = errors.New("campaign is paused or archived")
	ErrInvalidEventType          = errors.New("invalid event type")
	ErrEventNotFound             = errors.New("event not found")
	ErrNoActiveCampaignAvailable = errors.New("no active campaigns available")
	ErrNoBannersAvailable        = errors.New("no banners available for display")
	ErrInvalidBannerID           = errors.New("invalid banner ID format")
	ErrInvalidCampaignID         = errors.New("invalid campaign ID format")
	ErrInvalidDateRange          = errors.New("invalid date range: 'from' date must be before 'to' date")
	ErrBothDatesRequired         = errors.New("both 'from' and 'to' dates must be provided together")
)
