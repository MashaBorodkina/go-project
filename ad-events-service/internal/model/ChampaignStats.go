package model

import (
	"github.com/google/uuid"
)

type CampaignStats struct {
	CampaignID  uuid.UUID `json:"campaign_id"`
	Impressions int       `json:"impressions"`
	CampaignName string    `json:"campaign_name"`
	Budget      float64   `json:"budget"`
	Clicks      int       `json:"clicks"`
	CTR         float64   `json:"ctr"` // Click-Through Rate	
}
