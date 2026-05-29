package handler

import (
	"ad-events-service/internal/apperrors"
	"ad-events-service/internal/service"
	"errors"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventHandler struct {
	EventService *service.EventService
}

func NewEventHandler(eventService *service.EventService) *EventHandler {
	return &EventHandler{EventService: eventService}
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()
	var req struct {
		BannerID string `json:"banner_id" binding:"required,uuid"`
		Type     string `json:"type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	err := h.EventService.TrackEvent(c.Request.Context(), req.BannerID, req.Type, ip, userAgent)
	if err != nil {
		if errors.Is(err, apperrors.ErrBannerInactive) {
			Error(c, http.StatusBadRequest, "Banner is inactive")
			return
		}
		if errors.Is(err, apperrors.ErrCampaignNotFound) {
			Error(c, http.StatusNotFound, "Campaign not found")
			return
		}
		if errors.Is(err, apperrors.ErrCampaignInactive) {
			Error(c, http.StatusBadRequest, "Campaign is inactive")
			return
		}
		if errors.Is(err, apperrors.ErrInvalidEventType) {
			Error(c, http.StatusBadRequest, "Invalid event type")
			return
		}
		Error(c, http.StatusInternalServerError, "Failed to track event")
		return
	}
	Success(c, http.StatusCreated, "Event tracked successfully")
}

func (h *EventHandler) GetEventByID(c *gin.Context) {
	eventID := c.Param("event_id")
	_, err := uuid.Parse(eventID)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid banner ID format")
		return
	}
	events, err := h.EventService.GetEventByID(c.Request.Context(), eventID)
	if err != nil {
		Error(c, http.StatusInternalServerError, "Failed to retrieve events")
		return
	}
	Success(c, http.StatusOK, events)
}

func (h *EventHandler) GetEventsByBannerID(c *gin.Context) {
	bannerID := c.Param("banner_id")
	_, err := uuid.Parse(bannerID)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid banner ID format")
		return
	}
	eventType := c.Query("type")
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid limit value")
		return
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid offset value")
		return
	}

	events, err := h.EventService.GetEventsByBannerID(c.Request.Context(), bannerID, eventType, limitInt, offsetInt)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidEventType) {
			Error(c, http.StatusBadRequest, "Invalid event type")
			return
		}
		Error(c, http.StatusInternalServerError, "Failed to retrieve events")
		return
	}
	Success(c, http.StatusOK, events)
}

func (h *EventHandler) GetEventsByCampaignID(c *gin.Context) {
	campaignID := c.Param("campaign_id")
	_, err := uuid.Parse(campaignID)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid campaign ID format")
		return
	}
	eventType := c.Query("type")
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid limit value")
		return
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid offset value")
		return
	}

	events, err := h.EventService.GetEventsByCampaignID(c.Request.Context(), campaignID, eventType, limitInt, offsetInt)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidEventType) {
			Error(c, http.StatusBadRequest, "Invalid event type")
			return
		}
		Error(c, http.StatusInternalServerError, "Failed to retrieve events")
		return
	}
	Success(c, http.StatusOK, events)
}
