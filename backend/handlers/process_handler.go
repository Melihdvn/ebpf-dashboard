package handlers

import (
	"ebpf-dashboard/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProcessHandler struct {
	service services.ProcessService
}

func NewProcessHandler(service services.ProcessService) *ProcessHandler {
	return &ProcessHandler{service: service}
}

func (h *ProcessHandler) GetRecentProcesses(c *gin.Context) {
	// Get limit from query parameter, default to 100, max 1000
	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	processes, err := h.service.GetRecentProcesses(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(processes),
		"data":  processes,
	})
}
