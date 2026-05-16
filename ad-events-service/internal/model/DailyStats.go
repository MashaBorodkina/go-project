package model

import "time"

type DailyStats struct {
	Date        time.Time `json:"date"`
	Impressions int       `json:"impressions"`
	Clicks      int       `json:"clicks"`
	CTR         float64   `json:"ctr"` // Click-Through Rate
}