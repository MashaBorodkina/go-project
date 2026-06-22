package repository

import (
	"ad-events-service/internal/apperrors"
	"ad-events-service/internal/model"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetCampaignByID(ctx context.Context, id string) (*model.Campaign, error) {
	var campaign model.Campaign
	query := "SELECT id, name, budget, status, created_at, updated_at FROM campaigns WHERE id = $1"
	err := r.db.QueryRow(ctx, query, id).Scan(&campaign.ID, &campaign.Name, &campaign.Budget, &campaign.Status, &campaign.CreatedAt, &campaign.UpdatedAt)
	if err != nil {
		return nil, apperrors.ErrCampaignNotFound
	}
	return &campaign, nil
}

func (r *Repository) GetCampaignByName(ctx context.Context, name string) (*model.Campaign, error) {
	var campaign model.Campaign
	query := `SELECT id, name, budget, status, created_at, updated_at 
	FROM campaigns 
	WHERE name = $1`
	err := r.db.QueryRow(ctx, query, name).Scan(&campaign.ID, &campaign.Name, &campaign.Budget, &campaign.Status, &campaign.CreatedAt, &campaign.UpdatedAt)
	if err != nil {
		return nil, apperrors.ErrCampaignNotFound
	}
	return &campaign, nil
}

func (r *Repository) GetAllCampaigns(ctx context.Context) ([]*model.Campaign, error) {
	query := "SELECT id, name, budget, status, created_at, updated_at FROM campaigns"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all campaigns: %w", err)
	}
	defer rows.Close()

	var campaigns []*model.Campaign
	for rows.Next() {
		var campaign model.Campaign
		err := rows.Scan(&campaign.ID, &campaign.Name, &campaign.Budget, &campaign.Status, &campaign.CreatedAt, &campaign.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan campaign: %w", err)
		}
		campaigns = append(campaigns, &campaign)
	}

	return campaigns, nil
}

func (r *Repository) CreateCampaign(ctx context.Context, campaign *model.Campaign) error {
	query := `INSERT INTO campaigns (name, budget) 
	VALUES ($1, $2) 
	RETURNING id, name, budget, status, created_at`
	err := r.db.QueryRow(ctx, query, campaign.Name, campaign.Budget).Scan(&campaign.ID, &campaign.Name, &campaign.Budget, &campaign.Status, &campaign.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create campaign: %w", err)
	}
	return nil
}

func (r *Repository) UpdateCampaign(ctx context.Context, campaign *model.Campaign) error {
	query := "UPDATE campaigns SET name = $1, budget = $2, status = $3, updated_at = NOW() WHERE id = $4"
	_, err := r.db.Exec(ctx, query, campaign.Name, campaign.Budget, campaign.Status, campaign.ID)
	if err != nil {
		return fmt.Errorf("failed to update campaign: %w", err)
	}
	return nil
}

func (r *Repository) DeleteCampaign(ctx context.Context, id string) error {
	query := "DELETE FROM campaigns WHERE id = $1"
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete campaign: %w", err)
	}
	return nil
}
