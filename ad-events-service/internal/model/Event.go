package model

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID       `json:"id"`
	Banner_ID      uuid.UUID       `json:"banner_id"`
	Type      string    `json:"type"` // e.g., "impression", "click"
	Ip 	   string    `json:"ip"`
	User_Agent string    `json:"user_agent"`
	Timestamp time.Time `json:"timestamp"`
}