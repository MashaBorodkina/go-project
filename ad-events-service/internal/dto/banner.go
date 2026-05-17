package dto

type BannerPatchRequest struct {
	Title    *string `json:"title,omitempty"`
	ImageUrl *string `json:"image_url,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}