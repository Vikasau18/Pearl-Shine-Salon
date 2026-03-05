package handlers

import (
	"net/http"

	"saloon-backend/middleware"
	"saloon-backend/services"

	"github.com/gin-gonic/gin"
)

type PushHandler struct {
	Push *services.PushService
}

func NewPushHandler(push *services.PushService) *PushHandler {
	return &PushHandler{Push: push}
}

type subscribeRequest struct {
	Endpoint string `json:"endpoint" binding:"required"`
	Keys     struct {
		P256dh string `json:"p256dh" binding:"required"`
		Auth   string `json:"auth" binding:"required"`
	} `json:"keys" binding:"required"`
}

// Subscribe saves a push subscription for the authenticated user
func (h *PushHandler) Subscribe(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req subscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription data"})
		return
	}

	err := h.Push.SaveSubscription(c.Request.Context(), userID, req.Endpoint, req.Keys.P256dh, req.Keys.Auth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscribed to push notifications"})
}

type unsubscribeRequest struct {
	Endpoint string `json:"endpoint" binding:"required"`
}

// Unsubscribe removes a push subscription
func (h *PushHandler) Unsubscribe(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req unsubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.Push.RemoveSubscription(c.Request.Context(), userID, req.Endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed from push notifications"})
}

// GetVAPIDKey returns the public VAPID key for the frontend
func (h *PushHandler) GetVAPIDKey(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"public_key": h.Push.VAPIDPublicKey})
}
