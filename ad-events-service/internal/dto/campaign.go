package dto

import "github.com/google/uuid"

type PatchCampaignRequest struct {
	Name   *string  `json:"name,omitempty"`
	Budget *float64 `json:"budget,omitempty"`
	Status *string  `json:"status,omitempty"`
}

type CreateCampaignResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name" binding:"required"`
	Budget    float64   `json:"budget" binding:"required,gt=0"`
	Status    string    `json:"status" binding:"required,oneof=active paused archived"`
	CreatedAt string    `json:"created_at,omitempty"`
}
