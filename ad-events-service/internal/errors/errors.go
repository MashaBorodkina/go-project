package errors

import "errors"

var (
	ErrBannerNotFound =errors.New("banner not found")
	ErrCampaignNotFound =errors.New("campaign not found")
	ErrBannerInactive =errors.New("banner is inactive")
	ErrCampaignInactive =errors.New("campaign is paused or archived")
	ErrInvalidEventType =errors.New("invalid event type")
	ErrEventNotFound =errors.New("event not found")
)