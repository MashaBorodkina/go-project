package dto

import "time"

type TrackEventResponse struct {
	BannerID  string    `json:"banner_id"`
	CreatedAt time.Time `json:"recorded_at"`
}
