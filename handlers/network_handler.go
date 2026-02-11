package handlers

import (
	"ebpf-dashboard/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NetworkHandler struct {
	service services.NetworkService
}

func NewNetworkHandler(service services.NetworkService) *NetworkHandler {
	return &NetworkHandler{service: service}
}

func (h *NetworkHandler) GetRecentConnections(c *gin.Context) {
	// Get limit from query parameter, default to 50
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	connections, err := h.service.GetRecentConnections(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(connections),
		"data":  connections,
	})
}
