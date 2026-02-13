package handlers

import (
	"ebpf-dashboard/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TCPLifeHandler struct {
	service services.TCPLifeService
}

func NewTCPLifeHandler(service services.TCPLifeService) *TCPLifeHandler {
	return &TCPLifeHandler{service: service}
}

// GetTCPLifeEvents handles GET /api/metrics/tcplife
func (h *TCPLifeHandler) GetTCPLifeEvents(c *gin.Context) {
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

	events, err := h.service.GetRecentEvents(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(events),
		"data":  events,
	})
}
