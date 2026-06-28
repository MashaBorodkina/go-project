package dto

type PeriodeResponse struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type DailyStarsItemResponse struct {
	Date        string  `json:"date"`
	Impressions int     `json:"impressions"`
	Clicks      int     `json:"clicks"`
	CTR         float64 `json:"ctr"`
}

type DailyStatsResponse struct {
	CampaignID string                   `json:"campaign_id"`
	Period     PeriodeResponse          `json:"period"`
	Daily      []DailyStarsItemResponse `json:"daily"`
}
