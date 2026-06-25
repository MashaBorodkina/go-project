package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"ad-events-service/internal/apperrors"
	"ad-events-service/internal/service"
)

type AdHandler struct {
	AdService *service.AdService
}

func NewAdHandler(adService *service.AdService) *AdHandler {
	return &AdHandler{AdService: adService}
}

func (h *AdHandler) GetBannerForDisplay(c *gin.Context) {
	banner, err := h.AdService.GetBannerForDisplay(c.Request.Context())
	if err != nil {
		if errors.Is(err, apperrors.ErrNoActiveCampaignAvailable) {
			Error(c, http.StatusNotFound, "No active campaigns available")

			return
		}
		if errors.Is(err, apperrors.ErrNoBannersAvailable) {
			Error(c, http.StatusNotFound, "No banners available for display")

			return
		}
		Error(c, http.StatusInternalServerError, "Failed to get banner for display")

		return
	}
	Success(c, http.StatusOK, banner)
}
