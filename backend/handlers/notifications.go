package handlers

import (
	"context"
	"net/http"

	"saloon-backend/middleware"
	"saloon-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationHandler struct {
	DB *pgxpool.Pool
}

func NewNotificationHandler(db *pgxpool.Pool) *NotificationHandler {
	return &NotificationHandler{DB: db}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := middleware.GetUserID(c)

	rows, err := h.DB.Query(context.Background(),
		`SELECT id, user_id, type, title, message, is_read, appointment_id, created_at
		 FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT 50`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.IsRead, &n.AppointmentID, &n.CreatedAt)
		notifications = append(notifications, n)
	}
	if notifications == nil {
		notifications = []models.Notification{}
	}

	c.JSON(http.StatusOK, notifications)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	notificationID := c.Param("id")
	userID := middleware.GetUserID(c)

	_, err := h.DB.Exec(context.Background(),
		"UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2",
		notificationID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID := middleware.GetUserID(c)

	_, err := h.DB.Exec(context.Background(),
		"UPDATE notifications SET is_read = true WHERE user_id = $1 AND is_read = false", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var count int
	h.DB.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false", userID).Scan(&count)

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}
