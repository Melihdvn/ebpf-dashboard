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
	// Get limit from query parameter, default to 50
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
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
