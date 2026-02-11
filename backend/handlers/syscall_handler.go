package handlers

import (
	"ebpf-dashboard/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SyscallHandler struct {
	service *services.SyscallService
}

func NewSyscallHandler(service *services.SyscallService) *SyscallHandler {
	return &SyscallHandler{service: service}
}

// GetSyscallStats handles GET /api/metrics/syscalls
func (h *SyscallHandler) GetSyscallStats(c *gin.Context) {
	// Get limit from query params (default: 50)
	limitStr := c.Query("limit")
	limit := 50
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	stats, err := h.service.GetRecentStats(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
