package dto

type PatchCampaignRequest struct {
	Name   *string  `json:"name,omitempty"`
	Budget *float64 `json:"budget,omitempty"`
	Status *string  `json:"status,omitempty"`
}
