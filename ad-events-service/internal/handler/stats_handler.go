package handler

import (
	"ad-events-service/internal/apperrors"
	"ad-events-service/internal/service"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StatsHandler struct {
	StatsService *service.StatsService
}

func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{StatsService: statsService}
}

func (h *StatsHandler) GetBannerStatsByID(c *gin.Context) {
	bannerID := c.Param("id")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil && fromStr != "" {
		Error(c, http.StatusBadRequest, "Invalid 'from' date format. Use YYYY-MM-DD format.")
		return
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil && toStr != "" {
		Error(c, http.StatusBadRequest, "Invalid 'to' date format. Use YYYY-MM-DD format.")
		return
	}

	if _, err := uuid.Parse(bannerID); err != nil {
		Error(c, http.StatusBadRequest, "Invalid banner ID format")
		return
	}
	bannerStats, err := h.StatsService.GetBannerStatsByID(c.Request.Context(), bannerID, from, to)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidDateRange) || errors.Is(err, apperrors.ErrBothDatesRequired) {
			Error(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, apperrors.ErrBannerNotFound) {
			Error(c, http.StatusNotFound, "Banner not found")
			return
		}
		Error(c, http.StatusInternalServerError, "Failed to retrieve banner stats")
		return
	}
	Success(c, http.StatusOK, bannerStats)
}

func (h *StatsHandler) GetCampaignStatsByID(c *gin.Context) {
	campaignID := c.Param("id")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var from, to time.Time
	var err error

	if fromStr != "" {
		if from, err = time.Parse("2006-01-02", fromStr); err != nil {
			Error(c, http.StatusBadRequest, "Invalid 'from' date format. Use YYYY-MM-DD format.")
			return
		}
	}
	if toStr != "" {
		if to, err = time.Parse("2006-01-02", toStr); err != nil {
			Error(c, http.StatusBadRequest, "Invalid 'to' date format. Use YYYY-MM-DD format.")
			return
		}
	}

	if _, err := uuid.Parse(campaignID); err != nil {
		Error(c, http.StatusBadRequest, "Invalid campaign ID format")
		return
	}
	campaignStats, err := h.StatsService.GetCampaignStatsByID(c.Request.Context(), campaignID, from, to)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidDateRange) || errors.Is(err, apperrors.ErrBothDatesRequired) {
			Error(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, apperrors.ErrCampaignNotFound) {
			Error(c, http.StatusNotFound, "Campaign not found")
			return
		}
		Error(c, http.StatusInternalServerError, "Failed to retrieve campaign stats")
		return
	}
	Success(c, http.StatusOK, campaignStats)
}

func (h *StatsHandler) GetDailyStats(c *gin.Context) {
	campaignID := c.Param("id")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var from, to time.Time
	var err error

	if fromStr != "" {
		if from, err = time.Parse("2006-01-02", fromStr); err != nil {
			Error(c, http.StatusBadRequest, "Invalid 'from' date format. Use YYYY-MM-DD format.")
			return
		}
	}
	if toStr != "" {
		if to, err = time.Parse("2006-01-02", toStr); err != nil {
			Error(c, http.StatusBadRequest, "Invalid 'to' date format. Use YYYY-MM-DD format.")
			return
		}
	}

	if _, err := uuid.Parse(campaignID); err != nil {
		Error(c, http.StatusBadRequest, "Invalid campaign ID format")
		return
	}
	dailyStats, err := h.StatsService.GetDailyStats(c.Request.Context(), campaignID, from, to)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidDateRange) || errors.Is(err, apperrors.ErrBothDatesRequired) {
			Error(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, apperrors.ErrCampaignNotFound) {
			Error(c, http.StatusNotFound, "Campaign not found")
			return
		}
		Error(c, http.StatusInternalServerError, "Failed to retrieve daily stats")
		return
	}
	Success(c, http.StatusOK, dailyStats)
}
