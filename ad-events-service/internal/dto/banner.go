package dto

import (
	"time"
)

type BannerPatchRequest struct {
	Title    *string `json:"title,omitempty"`
	ImageUrl *string `json:"image_url,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type BannerResponseCreate struct {
	ID         string    `json:"id"`
	CampaignID string    `json:"campaign_id"`
	Title      string    `json:"title"`
	ImageUrl   string    `json:"image_url"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}
