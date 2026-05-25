package model

import (
	"github.com/google/uuid"
)

type BannerStats struct {
	BannerID    uuid.UUID `json:"banner_id"`
	Impressions int       `json:"impressions"`
	Clicks      int       `json:"clicks"`
	CTR         float64   `json:"ctr"` // Click-Through Rate
}
