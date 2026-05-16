package repository

import (
	"ad-events-service/internal/model"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) CreateEvent(ctx context.Context, event *model.Event) error {
	query := "INSERT INTO events (banner_id, type, user_agent, ip) VALUES ($1, $2, $3, $4) RETURNING id"
	err := r.db.QueryRow(ctx, query,  event.Banner_ID, event.Type, event.User_Agent, event.Ip).Scan(&event.ID)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	return nil
}

func (r *EventRepository) GetEventByID(ctx context.Context, ID string) (*model.Event, error) {
	var event model.Event
	query := "SELECT id, banner_id, type, created_at, ip, user_agent FROM events WHERE id = $1"
	err := r.db.QueryRow(ctx, query, ID).Scan(&event.ID, &event.Banner_ID, &event.Type, &event.CreatedAt, &event.Ip, &event.User_Agent)
	if err != nil {
		return nil, fmt.Errorf("failed to get event by ID: %w", err)
	}
	return &event, nil	
}

func (r *EventRepository) GetAllEvents(ctx context.Context) ([]*model.Event, error) {
	query := "SELECT id, banner_id, type, created_at, ip, user_agent FROM events"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all events: %w", err)
	}
	defer rows.Close()
	var events []*model.Event
	for rows.Next() {
		var event model.Event
		err := rows.Scan(&event.ID, &event.Banner_ID, &event.Type, &event.CreatedAt, &event.Ip, &event.User_Agent)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over events: %w", err)
	}
	return events, nil
}




