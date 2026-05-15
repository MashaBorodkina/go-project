package model

import (
	"time"

	"github.com/google/uuid"
)

type Campaign struct {
	ID 		uuid.UUID   `json:"id"`
	Name 	string `json:"name"`
	Budget 	float64 `json:"budget"`
	Status 	string    `json:"status"` // active, paused, archived
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}