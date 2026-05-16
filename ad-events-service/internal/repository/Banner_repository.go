package repository

import (
	"ad-events-service/internal/model"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BannerRepository struct {
	db *pgxpool.Pool
}

func NewBannerRepository(db *pgxpool.Pool) *BannerRepository {
	return &BannerRepository{db: db}
}

func (r *BannerRepository) GetBannerByID(ctx context.Context, id string) (*model.Banner, error) {
	var banner model.Banner
	query := "SELECT id, campaign_id, name, image_url, created_at, updated_at FROM banners WHERE id = $1"
	err := r.db.QueryRow(ctx, query, id).Scan(&banner.ID, &banner.CampaignID, &banner.Title, &banner.ImageUrl, &banner.CreatedAt, &banner.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get banner by ID: %w", err)
	}
	return &banner, nil
}

func (r *BannerRepository) GetAllBanners(ctx context.Context) ([]*model.Banner, error) {
	query := "SELECT id, campaign_id, name, image_url, created_at, updated_at FROM banners"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all banners: %w", err)
	}
	defer rows.Close()

	var banners []*model.Banner
	for rows.Next() {
		var banner model.Banner
		err := rows.Scan(&banner.ID, &banner.CampaignID, &banner.Title, &banner.ImageUrl, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan banner: %w", err)
		}
		banners = append(banners, &banner)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over banners: %w", err)
	}

	return banners, nil
}

func (r *BannerRepository) CreateBanner(ctx context.Context, banner *model.Banner) error {
	query := "INSERT INTO banners (campaign_id, title, image_url) VALUES ($1, $2, $3, $4) RETURNING id"
	err := r.db.QueryRow(ctx, query, banner.CampaignID, banner.Title, banner.ImageUrl, banner.IsActive).Scan(&banner.ID)	
	if err != nil {
		return fmt.Errorf("failed to create banner: %w", err)
	}
	return nil
}

func (r *BannerRepository) UpdateBanner (ctx context.Context, banner *model.Banner) error {
	quey := "UPDATE banners SET campaign_id = $1, title = $2, image_url = $3, is_active = $4, updated_at = NOW() WHERE id = $5"
	_, err := r.db.Exec(ctx, quey, banner.CampaignID, banner.Title, banner.ImageUrl, banner.IsActive, banner.ID)
	if err != nil {
		return fmt.Errorf("failed to update banner: %w", err)
	}
	return nil	
}

func (r *BannerRepository) DeleteBanner(ctx context.Context, id string) error {
	query := "DELETE FROM banners WHERE id = $1"
	_, err := r.db.Exec(ctx, query, id)		
	if err != nil {
		return fmt.Errorf("failed to delete banner: %w", err)
	}
	return nil
}

func (r *BannerRepository) UpdateBannerPartial (ctx context.Context, banner *model.Banner) error {
	query := "UPDATE banners SET campaign_id = COALESCE($1, campaign_id), title = COALESCE($2, title), image_url = COALESCE($3, image_url), is_active = COALESCE($4, is_active), updated_at = NOW() WHERE id = $5"
	_, err := r.db.Exec(ctx, query, banner.CampaignID, banner.Title, banner.ImageUrl, banner.IsActive, banner.ID)
	if err != nil {
		return fmt.Errorf("failed to update banner: %w", err)
	}
	return nil	
}