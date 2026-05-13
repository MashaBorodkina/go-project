package main

import (
	"time"

	"github.com/google/uuid"
)
	
	type Banner struct {
	ID        uuid.UUID       `json:"id"`
	CampaignID uuid.UUID       `json:"campaign_id"`
	Image_url string          `json:"image_url"`
	Title    string          `json:"title"`
	Is_active bool            `json:"is_active"`// true | false
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}