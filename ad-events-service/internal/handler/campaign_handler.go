package handler

import (
	"ad-events-service/internal/apperrors"
	"ad-events-service/internal/dto"
	"ad-events-service/internal/model"
	"ad-events-service/internal/service"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CampaignHandler struct {
	CampService *service.CampaignService
}

func NewCampaignHandler(campService *service.CampaignService) *CampaignHandler {
	return &CampaignHandler{CampService: campService}
}

func (h *CampaignHandler) GetCampaignByID(c *gin.Context) {
	campID := c.Param("id")
	_, err := uuid.Parse(campID)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid campaign ID format")
		return
	}
	camp, err := h.CampService.GetCampaignByID(c.Request.Context(), campID)
	if err != nil {
		if errors.Is(err, apperrors.ErrCampaignNotFound) {
			Error(c, http.StatusNotFound, "Campaign not found")
			return
		}
		Error(c, http.StatusInternalServerError, "Failed to retrieve campaign")
		return
	}
	Success(c, http.StatusOK, camp)
}

func (h *CampaignHandler) GetCampaignByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		Error(c, http.StatusBadRequest, "Campaign name query parameter is required")
		return
	}
	camp, err := h.CampService.GetCampaignByName(c.Request.Context(), name)
	if err != nil {
		if errors.Is(err, apperrors.ErrCampaignNotFound) {
			Error(c, http.StatusNotFound, "Campaign not found")
			return
		}
		Error(c, http.StatusInternalServerError, "Failed to retrieve campaign")
		return
	}
	Success(c, http.StatusOK, camp)
}

func (h *CampaignHandler) GetAllCampaigns(c *gin.Context) {
	camps, err := h.CampService.GetAllCampaigns(c.Request.Context())
	if err != nil {
		Error(c, http.StatusInternalServerError, "Failed to retrieve campaigns")
		return
	}
	Success(c, http.StatusOK, camps)
}

func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var req struct {
		Name   string  `json:"name" binding:"required"`
		Budget float64 `json:"budget" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	camp := &model.Campaign{
		Name:   req.Name,
		Budget: req.Budget,
	}
	err := h.CampService.CreateCampaign(c.Request.Context(), camp)
	if err != nil {
		fmt.Printf("Error creating campaign: %v\n", err)
		Error(c, http.StatusInternalServerError, "Failed to create campaign")
		return
	}
	response := dto.CreateCampaignResponse{
		ID:        camp.ID,
		Name:      camp.Name,
		Budget:    camp.Budget,
		Status:    camp.Status,
		CreatedAt: camp.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	Success(c, http.StatusCreated, response)
}

func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	campID := c.Param("id")
	parsID, err := uuid.Parse(campID)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid campaign ID format")
		return
	}
	var req struct {
		Name   string  `json:"name" binding:"required"`
		Budget float64 `json:"budget" binding:"required,gt=0"`
		Status string  `json:"status" binding:"required,oneof=active paused archived"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	camp := &model.Campaign{
		ID:     parsID,
		Name:   req.Name,
		Budget: req.Budget,
		Status: req.Status,
	}
	err = h.CampService.UpdateCampaign(c.Request.Context(), camp)
	if err != nil {
		if errors.Is(err, apperrors.ErrCampaignNotFound) {
			Error(c, http.StatusNotFound, "Campaign not found")
			return
		}
		Error(c, http.StatusInternalServerError, "Failed to update campaign")
		return
	}
	Success(c, http.StatusOK, camp)
}

func (h *CampaignHandler) DeleteCampaign(c *gin.Context) {
	campID := c.Param("id")
	_, err := uuid.Parse(campID)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid campaign ID format")
		return
	}
	err = h.CampService.DeleteCampaign(c.Request.Context(), campID)
	if err != nil {
		if errors.Is(err, apperrors.ErrCampaignNotFound) {
			Error(c, http.StatusNotFound, "Campaign not found")
			return
		}
		Error(c, http.StatusInternalServerError, "Failed to delete campaign")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *CampaignHandler) PatchCampaign(c *gin.Context) {
	campID := c.Param("id")
	_, err := uuid.Parse(campID)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid campaign ID format")
		return
	}
	var req dto.PatchCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	updatedCamp, err := h.CampService.PatchCampaign(c.Request.Context(), campID, &req)
	if err != nil {
		if errors.Is(err, apperrors.ErrCampaignNotFound) {
			Error(c, http.StatusNotFound, "Campaign not found")
			return
		}
		Error(c, http.StatusInternalServerError, "Failed to update campaign")
		return
	}
	Success(c, http.StatusOK, updatedCamp)
}
