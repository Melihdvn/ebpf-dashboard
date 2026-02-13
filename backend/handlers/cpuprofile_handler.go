package handlers

import (
	"ebpf-dashboard/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CPUProfileHandler struct {
	service services.CPUProfileService
}

func NewCPUProfileHandler(service services.CPUProfileService) *CPUProfileHandler {
	return &CPUProfileHandler{service: service}
}

// GetCPUProfiles handles GET /api/metrics/cpuprofile
func (h *CPUProfileHandler) GetCPUProfiles(c *gin.Context) {
	// Get limit from query params (default: 100, max: 1000)
	limitStr := c.Query("limit")
	limit := 100
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			if parsedLimit > 1000 {
				limit = 1000
			} else {
				limit = parsedLimit
			}
		}
	}

	profiles, err := h.service.GetRecentProfiles(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(profiles),
		"data":  profiles,
	})
}
