package model

import (
	"github.com/google/uuid"
)

type CampaignStats struct {
	CampaignID   uuid.UUID `json:"campaign_id"`
	CampaignName string    `json:"campaign_name"`
	Impressions  int       `json:"impressions"`
	Clicks       int       `json:"clicks"`
	CTR          float64   `json:"ctr"` // Click-Through Rate
	Budget       float64   `json:"budget"`
}
