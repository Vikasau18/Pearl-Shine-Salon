package handlers

import (
	"context"
	"net/http"

	"saloon-backend/middleware"
	"saloon-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WaitlistHandler struct {
	DB *pgxpool.Pool
}

func NewWaitlistHandler(db *pgxpool.Pool) *WaitlistHandler {
	return &WaitlistHandler{DB: db}
}

func (h *WaitlistHandler) JoinWaitlist(c *gin.Context) {
	var req models.JoinWaitlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerID := middleware.GetUserID(c)

	var staffID *string
	if req.StaffID != "" {
		staffID = &req.StaffID
	}
	var preferredTime *string
	if req.PreferredTime != "" {
		preferredTime = &req.PreferredTime
	}

	var w models.Waitlist
	err := h.DB.QueryRow(context.Background(),
		`INSERT INTO waitlist (customer_id, salon_id, service_id, staff_id, preferred_date, preferred_time)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, customer_id, salon_id, service_id, staff_id, preferred_date::text, preferred_time::text, status, created_at`,
		customerID, req.SalonID, req.ServiceID, staffID, req.PreferredDate, preferredTime,
	).Scan(&w.ID, &w.CustomerID, &w.SalonID, &w.ServiceID, &w.StaffID,
		&w.PreferredDate, &w.PreferredTime, &w.Status, &w.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join waitlist", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, w)
}

func (h *WaitlistHandler) LeaveWaitlist(c *gin.Context) {
	waitlistID := c.Param("id")
	userID := middleware.GetUserID(c)

	result, err := h.DB.Exec(context.Background(),
		"DELETE FROM waitlist WHERE id = $1 AND customer_id = $2", waitlistID, userID)
	if err != nil || result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Waitlist entry not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from waitlist"})
}

func (h *WaitlistHandler) GetSalonWaitlist(c *gin.Context) {
	salonID := c.Param("salon_id")

	rows, err := h.DB.Query(context.Background(),
		`SELECT w.id, w.customer_id, w.salon_id, w.service_id, w.staff_id, 
		 w.preferred_date::text, w.preferred_time::text, w.status, w.created_at,
		 u.name, sv.name
		 FROM waitlist w
		 JOIN users u ON u.id = w.customer_id
		 JOIN services sv ON sv.id = w.service_id
		 WHERE w.salon_id = $1 AND w.status = 'waiting'
		 ORDER BY w.created_at`, salonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch waitlist"})
		return
	}
	defer rows.Close()

	type WaitlistEntry struct {
		models.Waitlist
		CustomerName string `json:"customer_name"`
		ServiceName  string `json:"service_name"`
	}

	var entries []WaitlistEntry
	for rows.Next() {
		var e WaitlistEntry
		rows.Scan(&e.ID, &e.CustomerID, &e.SalonID, &e.ServiceID, &e.StaffID,
			&e.PreferredDate, &e.PreferredTime, &e.Status, &e.CreatedAt,
			&e.CustomerName, &e.ServiceName)
		entries = append(entries, e)
	}
	if entries == nil {
		entries = []WaitlistEntry{}
	}

	c.JSON(http.StatusOK, entries)
}
