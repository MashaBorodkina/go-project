package handler

import (
	"ad-events-service/internal/apperrors"
	"ad-events-service/internal/dto"
	"ad-events-service/internal/model"
	"ad-events-service/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BannerHandler struct {
	BanService *service.BannerService
}

func NewBannerHandler(banService *service.BannerService) *BannerHandler {
	return &BannerHandler{BanService: banService}
}

func (h *BannerHandler) GetBannerByID(c *gin.Context) {
	banID := c.Param("id")
	_, err := uuid.Parse(banID)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid banner ID format")
		return
	}

	ban, err := h.BanService.GetBannerByID(c.Request.Context(), banID)
	if err != nil {
		if errors.Is(err, apperrors.ErrBannerNotFound) {
			Error(c, http.StatusNotFound, "Banner not found")
			return
		}
		Error(c, http.StatusInternalServerError, "Failed to retrieve banner")
		return
	}

	Success(c, http.StatusOK, ban)
}

func (h *BannerHandler) GetAllBanners(c *gin.Context) {
	bans, err := h.BanService.GetAllBanners(c.Request.Context())
	if err != nil {
		Error(c, http.StatusInternalServerError, "Failed to retrieve banners")
		return
	}
	Success(c, http.StatusOK, bans)
}

func (h *BannerHandler) CreateBanner(c *gin.Context) {
	campaignID := c.Param("campaign_id")
	parsId, err := uuid.Parse(campaignID)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid campaign ID format")
		return
	}
	var req struct {
		Title    string `json:"title" binding:"required"`
		ImageUrl string `json:"image_url" binding:"required"`
		IsActive bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	banner := &model.Banner{
		CampaignID: parsId,
		Title:      req.Title,
		ImageUrl:   req.ImageUrl,
		IsActive:   req.IsActive,
	}
	err = h.BanService.CreateBanner(c.Request.Context(), banner)
	if err != nil {
		Error(c, http.StatusInternalServerError, "Failed to create banner")
		return
	}
	Success(c, http.StatusCreated, banner)
}

func (h *BannerHandler) UpdateBanner(c *gin.Context) {
	banID := c.Param("id")
	_, err := uuid.Parse(banID)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid banner ID format")
		return
	}
	var req struct {
		CampaignID string `json:"campaign_id" binding:"required"`
		Title      string `json:"title" binding:"required"`
		ImageUrl   string `json:"image_url" binding:"required"`
		IsActive   bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	parsId, err := uuid.Parse(req.CampaignID)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid campaign ID format")
		return
	}

	banner := &model.Banner{
		ID:         uuid.MustParse(banID),
		CampaignID: parsId,
		IsActive:   req.IsActive,
		Title:      req.Title,
		ImageUrl:   req.ImageUrl,
	}
	err = h.BanService.UpdateBanner(c.Request.Context(), banner)
	if err != nil {
		if errors.Is(err, apperrors.ErrBannerNotFound) {
			Error(c, http.StatusNotFound, "Banner not found")
			return
		}
		Error(c, http.StatusInternalServerError, "Failed to update banner")
		return
	}
	Success(c, http.StatusOK, banner)
}

func (h *BannerHandler) DeleteBanner(c *gin.Context) {
	banID := c.Param("id")
	_, err := uuid.Parse(banID)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid banner ID format")
		return
	}
	err = h.BanService.DeleteBanner(c.Request.Context(), banID)
	if err != nil {
		if errors.Is(err, apperrors.ErrBannerNotFound) {
			Error(c, http.StatusNotFound, "Banner not found")
			return
		}
		Error(c, http.StatusInternalServerError, "Failed to delete banner")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *BannerHandler) PatchBanner(c *gin.Context) {
	banID := c.Param("id")
	_, err := uuid.Parse(banID)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid banner ID format")
		return
	}
	var req dto.BannerPatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	updatedBan, err := h.BanService.PatchBanner(c.Request.Context(), banID, &req)
	if err != nil {
		if errors.Is(err, apperrors.ErrBannerNotFound) {
			Error(c, http.StatusNotFound, "Banner not found")
			return
		}
		Error(c, http.StatusInternalServerError, "Failed to update banner")
		return
	}
	Success(c, http.StatusOK, updatedBan)
}
