package handlers

import (
	"ebpf-dashboard/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DiskHandler struct {
	service services.DiskService
}

func NewDiskHandler(service services.DiskService) *DiskHandler {
	return &DiskHandler{service: service}
}

func (h *DiskHandler) GetLatestLatency(c *gin.Context) {
	// Get limit from query parameter, default to 20 (histogram buckets)
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	latencies, err := h.service.GetLatestLatency(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(latencies),
		"data":  latencies,
	})
}
