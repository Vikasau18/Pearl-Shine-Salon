package handlers

import (
	"context"
	"fmt"
	"net/http"

	"saloon-backend/models"
	"saloon-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClosureHandler struct {
	DB   *pgxpool.Pool
	Push *services.PushService
}

func NewClosureHandler(db *pgxpool.Pool, push *services.PushService) *ClosureHandler {
	return &ClosureHandler{DB: db, Push: push}
}

// CreateClosure creates a new salon closure and returns any conflicting appointments
func (h *ClosureHandler) CreateClosure(c *gin.Context) {
	salonID := c.Param("salon_id")

	var req models.CreateClosureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert closure
	var closure models.SalonClosure
	err := h.DB.QueryRow(context.Background(),
		`INSERT INTO salon_closures (salon_id, start_date, end_date, start_time, end_time, reason)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, salon_id, start_date::text, end_date::text,
		 start_time::text, end_time::text, reason, created_at`,
		salonID, req.StartDate, req.EndDate, req.StartTime, req.EndTime, req.Reason,
	).Scan(&closure.ID, &closure.SalonID, &closure.StartDate, &closure.EndDate,
		&closure.StartTime, &closure.EndTime, &closure.Reason, &closure.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create closure", "details": err.Error()})
		return
	}

	// Find conflicting appointments
	conflictQuery := `
		SELECT a.id, u.name, sv.name, a.appointment_date::text, a.start_time::text, a.status
		FROM appointments a
		JOIN users u ON u.id = a.customer_id
		JOIN services sv ON sv.id = a.service_id
		WHERE a.salon_id = $1
		  AND a.appointment_date >= $2::date AND a.appointment_date <= $3::date
		  AND a.status NOT IN ('cancelled', 'completed', 'no_show')
	`
	rows, err := h.DB.Query(context.Background(), conflictQuery, salonID, req.StartDate, req.EndDate)

	var conflicts []models.ConflictingAppointment
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ca models.ConflictingAppointment
			rows.Scan(&ca.ID, &ca.CustomerName, &ca.ServiceName, &ca.AppointmentDate, &ca.StartTime, &ca.Status)
			conflicts = append(conflicts, ca)
		}
	}
	if conflicts == nil {
		conflicts = []models.ConflictingAppointment{}
	}

	c.JSON(http.StatusCreated, models.CreateClosureResponse{
		Closure:             closure,
		ConflictingCount:    len(conflicts),
		ConflictingBookings: conflicts,
	})
}

// GetClosures lists all closures for a salon
func (h *ClosureHandler) GetClosures(c *gin.Context) {
	salonID := c.Param("salon_id")

	rows, err := h.DB.Query(context.Background(),
		`SELECT id, salon_id, start_date::text, end_date::text,
		 start_time::text, end_time::text, reason, created_at
		 FROM salon_closures WHERE salon_id = $1
		 ORDER BY start_date ASC`, salonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch closures"})
		return
	}
	defer rows.Close()

	var closures []models.SalonClosure
	for rows.Next() {
		var cl models.SalonClosure
		rows.Scan(&cl.ID, &cl.SalonID, &cl.StartDate, &cl.EndDate,
			&cl.StartTime, &cl.EndTime, &cl.Reason, &cl.CreatedAt)
		closures = append(closures, cl)
	}
	if closures == nil {
		closures = []models.SalonClosure{}
	}

	c.JSON(http.StatusOK, closures)
}

// DeleteClosure removes a closure
func (h *ClosureHandler) DeleteClosure(c *gin.Context) {
	salonID := c.Param("salon_id")
	closureID := c.Param("closure_id")

	tag, err := h.DB.Exec(context.Background(),
		"DELETE FROM salon_closures WHERE id = $1 AND salon_id = $2", closureID, salonID)
	if err != nil || tag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Closure not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Closure deleted"})
}

// CancelConflicting bulk-cancels appointments that conflict with a closure
func (h *ClosureHandler) CancelConflicting(c *gin.Context) {
	salonID := c.Param("salon_id")
	closureID := c.Param("closure_id")

	// Get closure dates
	var startDate, endDate string
	err := h.DB.QueryRow(context.Background(),
		"SELECT start_date::text, end_date::text FROM salon_closures WHERE id = $1 AND salon_id = $2",
		closureID, salonID).Scan(&startDate, &endDate)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Closure not found"})
		return
	}

	// Cancel all conflicting appointments and get customer IDs for notification
	updateQuery := `
		UPDATE appointments SET status = 'cancelled'
		WHERE salon_id = $1
		  AND appointment_date >= $2::date AND appointment_date <= $3::date
		  AND status NOT IN ('cancelled', 'completed', 'no_show')
		RETURNING customer_id, appointment_date::text, start_time::text
	`
	rows, err := h.DB.Query(context.Background(), updateQuery, salonID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update appointments"})
		return
	}
	defer rows.Close()

	// Get salon name for the notification
	var salonName string
	_ = h.DB.QueryRow(context.Background(), "SELECT name FROM salons WHERE id = $1", salonID).Scan(&salonName)

	var cancelledCount int64
	for rows.Next() {
		var customerID, apptDate, apptTime string
		if err := rows.Scan(&customerID, &apptDate, &apptTime); err == nil {
			cancelledCount++
			// Send push notification
			if h.Push != nil {
				h.Push.SendToUser(context.Background(), customerID, services.PushPayload{
					Title: "Appointment Cancelled",
					Body:  fmt.Sprintf("%s is closed on %s %s. Your appointment has been cancelled.", salonName, apptDate, apptTime),
					URL:   "/appointments",
				})
			}
			// Also create in-app notification record
			_, _ = h.DB.Exec(context.Background(),
				`INSERT INTO notifications (user_id, title, message, type)
				 VALUES ($1, $2, $3, $4)`,
				customerID, "Salon Closure Notification",
				fmt.Sprintf("Due to an emergency closure, your appointment at %s on %s at %s was cancelled. Please reschedule at your convenience.", salonName, apptDate, apptTime),
				"appointment_cancelled")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Conflicting appointments cancelled and customers notified",
		"cancelled_count": cancelledCount,
	})
}
