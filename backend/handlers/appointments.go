package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"saloon-backend/middleware"
	"saloon-backend/models"
	"saloon-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppointmentHandler struct {
	DB              *pgxpool.Pool
	BookingService  *services.BookingService
	PushService     *services.PushService
	WaitlistService *services.WaitlistService
}

func NewAppointmentHandler(db *pgxpool.Pool, scheduler *services.Scheduler) *AppointmentHandler {
	return &AppointmentHandler{
		DB:              db,
		BookingService:  services.NewBookingService(db, scheduler),
		WaitlistService: services.NewWaitlistService(db, nil), // PushService set later
	}
}

func (h *AppointmentHandler) SetPushService(ps *services.PushService) {
	h.PushService = ps
	if h.WaitlistService != nil {
		h.WaitlistService.PushService = ps
	}
	if h.BookingService != nil {
		h.BookingService.SetPushService(ps)
	}
}

func (h *AppointmentHandler) BookAppointment(c *gin.Context) {
	var req models.BookAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerID := middleware.GetUserID(c)

	appt, payment, err := h.BookingService.BookAppointment(c.Request.Context(), customerID, req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"appointment": appt,
		"payment":     payment,
	})

	// Send immediate push notification (non-blocking)
	if h.PushService != nil {
		go h.PushService.SendToUser(context.Background(), customerID, services.PushPayload{
			Title: "Booking Confirmed! ✅",
			Body:  fmt.Sprintf("Your appointment for %s on %s at %s has been confirmed.", appt.ServiceName, appt.AppointmentDate, appt.StartTime),
			Icon:  "/vite.svg",
			URL:   "/appointments",
		})
	}
}

func (h *AppointmentHandler) GetMyAppointments(c *gin.Context) {
	userID := middleware.GetUserID(c)
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := `SELECT a.id, a.customer_id, a.salon_id, a.staff_id, a.service_id,
		a.appointment_date::text, a.start_time::text, a.end_time::text, a.status,
		COALESCE(a.notes,''), a.promo_code_id, a.created_at, a.updated_at,
		s.name, st.name, sv.name, sv.price, COALESCE(p.total, 0)
		FROM appointments a
		JOIN salons s ON s.id = a.salon_id
		JOIN staff st ON st.id = a.staff_id
		JOIN services sv ON sv.id = a.service_id
		LEFT JOIN payments p ON p.appointment_id = a.id
		WHERE a.customer_id = $1`

	args := []interface{}{userID}
	argIdx := 2

	if status == "upcoming" {
		query += " AND a.status IN ('confirmed', 'pending') AND (a.appointment_date > CURRENT_DATE OR (a.appointment_date = CURRENT_DATE AND a.end_time >= CURRENT_TIME))"
	} else if status != "" {
		query += " AND a.status = $" + strconv.Itoa(argIdx)
		args = append(args, status)
		argIdx++
	}

	query += " ORDER BY a.appointment_date DESC, a.start_time DESC"
	query += " LIMIT $" + strconv.Itoa(argIdx) + " OFFSET $" + strconv.Itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
		return
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var a models.Appointment
		rows.Scan(&a.ID, &a.CustomerID, &a.SalonID, &a.StaffID, &a.ServiceID,
			&a.AppointmentDate, &a.StartTime, &a.EndTime, &a.Status,
			&a.Notes, &a.PromoCodeID, &a.CreatedAt, &a.UpdatedAt,
			&a.SalonName, &a.StaffName, &a.ServiceName, &a.ServicePrice, &a.TotalPrice)
		appointments = append(appointments, a)
	}
	if appointments == nil {
		appointments = []models.Appointment{}
	}

	c.JSON(http.StatusOK, appointments)
}

func (h *AppointmentHandler) GetSalonAppointments(c *gin.Context) {
	salonID := c.Param("salon_id")
	staffID := c.Query("staff_id")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	status := c.Query("status")

	query := `SELECT a.id, a.customer_id, a.salon_id, a.staff_id, a.service_id,
		a.appointment_date::text, a.start_time::text, a.end_time::text, a.status,
		COALESCE(a.notes,''), a.promo_code_id, a.created_at, a.updated_at,
		s.name, st.name, sv.name, sv.price, u.name, u.email, COALESCE(u.phone,'')
		FROM appointments a
		JOIN salons s ON s.id = a.salon_id
		JOIN staff st ON st.id = a.staff_id
		JOIN services sv ON sv.id = a.service_id
		JOIN users u ON u.id = a.customer_id
		WHERE a.salon_id = $1`

	args := []interface{}{salonID}
	argIdx := 2

	if staffID != "" {
		query += " AND a.staff_id = $" + strconv.Itoa(argIdx)
		args = append(args, staffID)
		argIdx++
	}
	if dateFrom != "" {
		query += " AND a.appointment_date >= $" + strconv.Itoa(argIdx)
		args = append(args, dateFrom)
		argIdx++
	}
	if dateTo != "" {
		query += " AND a.appointment_date <= $" + strconv.Itoa(argIdx)
		args = append(args, dateTo)
		argIdx++
	}
	if status != "" {
		query += " AND a.status = $" + strconv.Itoa(argIdx)
		args = append(args, status)
		argIdx++
	}

	query += " ORDER BY a.appointment_date, a.start_time"

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
		return
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var a models.Appointment
		rows.Scan(&a.ID, &a.CustomerID, &a.SalonID, &a.StaffID, &a.ServiceID,
			&a.AppointmentDate, &a.StartTime, &a.EndTime, &a.Status,
			&a.Notes, &a.PromoCodeID, &a.CreatedAt, &a.UpdatedAt,
			&a.SalonName, &a.StaffName, &a.ServiceName, &a.ServicePrice,
			&a.CustomerName, &a.CustomerEmail, &a.CustomerPhone)
		appointments = append(appointments, a)
	}
	if appointments == nil {
		appointments = []models.Appointment{}
	}

	c.JSON(http.StatusOK, appointments)
}

func (h *AppointmentHandler) RescheduleAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	var req models.RescheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)

	// Get current appointment
	var staffID, serviceID, currentStatus string
	err := h.DB.QueryRow(context.Background(),
		"SELECT staff_id, service_id, status FROM appointments WHERE id = $1 AND customer_id = $2",
		appointmentID, userID).Scan(&staffID, &serviceID, &currentStatus)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	if currentStatus == "cancelled" || currentStatus == "completed" || currentStatus == "no_show" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot reschedule a " + currentStatus + " appointment"})
		return
	}

	// Get service duration
	var durationMinutes int
	h.DB.QueryRow(context.Background(),
		"SELECT duration_minutes FROM services WHERE id = $1", serviceID).Scan(&durationMinutes)

	// Calculate new end time
	startTime, _ := time.Parse("15:04", req.StartTime)
	endTime := startTime.Add(time.Duration(durationMinutes) * time.Minute)
	endTimeStr := endTime.Format("15:04")

	// Check for conflicts
	var conflictCount int
	h.DB.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM appointments 
		 WHERE staff_id = $1 AND appointment_date = $2 AND id != $3
		 AND status NOT IN ('cancelled', 'no_show')
		 AND (start_time, end_time) OVERLAPS ($4::time, $5::time)`,
		staffID, req.Date, appointmentID, req.StartTime, endTimeStr).Scan(&conflictCount)

	if conflictCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "New time slot conflicts with existing appointment"})
		return
	}

	_, err = h.DB.Exec(context.Background(),
		`UPDATE appointments SET appointment_date = $1, start_time = $2, end_time = $3, 
		 status = 'confirmed', updated_at = NOW() WHERE id = $4`,
		req.Date, req.StartTime, endTimeStr, appointmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reschedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment rescheduled successfully"})
}

func (h *AppointmentHandler) CancelAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)

	query := "UPDATE appointments SET status = 'cancelled', updated_at = NOW() WHERE id = $1"
	args := []interface{}{appointmentID}

	if userRole != "salon_owner" && userRole != "admin" {
		query += " AND customer_id = $2"
		args = append(args, userID)
	}

	result, err := h.DB.Exec(context.Background(), query, args...)
	if err != nil || result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found or already cancelled"})
		return
	}

	// Refund payment
	h.DB.Exec(context.Background(),
		"UPDATE payments SET status = 'refunded', updated_at = NOW() WHERE appointment_id = $1",
		appointmentID)

	c.JSON(http.StatusOK, gin.H{"message": "Appointment cancelled"})

	// Notify waitlist
	if h.WaitlistService != nil {
		var salonID, serviceID, staffID, dateStr, startTime string
		err := h.DB.QueryRow(context.Background(),
			"SELECT salon_id, service_id, staff_id, appointment_date::text, start_time::text FROM appointments WHERE id = $1",
			appointmentID).Scan(&salonID, &serviceID, &staffID, &dateStr, &startTime)
		if err == nil {
			go h.WaitlistService.AutoAssignWaitlist(context.Background(), salonID, serviceID, staffID, dateStr, startTime)
		}
	}
}

func (h *AppointmentHandler) MarkNoShow(c *gin.Context) {
	appointmentID := c.Param("id")

	_, err := h.DB.Exec(context.Background(),
		"UPDATE appointments SET status = 'no_show', updated_at = NOW() WHERE id = $1",
		appointmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark no-show"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Marked as no-show"})

	// Notify waitlist
	if h.WaitlistService != nil {
		var salonID, serviceID, staffID, dateStr, startTime string
		err := h.DB.QueryRow(context.Background(),
			"SELECT salon_id, service_id, staff_id, appointment_date::text, start_time::text FROM appointments WHERE id = $1",
			appointmentID).Scan(&salonID, &serviceID, &staffID, &dateStr, &startTime)
		if err == nil {
			go h.WaitlistService.AutoAssignWaitlist(context.Background(), salonID, serviceID, staffID, dateStr, startTime)
		}
	}
}

func (h *AppointmentHandler) ApproveAppointment(c *gin.Context) {
	appointmentID := c.Param("id")

	// 1. Update status to confirmed
	var customerID, serviceName, salonName, startTimeStr, apptDateStr string
	err := h.DB.QueryRow(context.Background(),
		`UPDATE appointments 
		 SET status = 'confirmed', updated_at = NOW() 
		 WHERE id = $1 AND status = 'pending'
		 RETURNING customer_id, 
		           (SELECT name FROM services WHERE id = appointments.service_id),
		           (SELECT name FROM salons WHERE id = appointments.salon_id),
		           start_time::text, appointment_date::text`,
		appointmentID).Scan(&customerID, &serviceName, &salonName, &startTimeStr, &apptDateStr)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found or already confirmed/cancelled"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment approved successfully"})

	// 2. Send immediate push notification
	if h.PushService != nil {
		go h.PushService.SendToUser(context.Background(), customerID, services.PushPayload{
			Title: "Booking Confirmed! ✅",
			Body:  fmt.Sprintf("Your appointment for %s at %s on %s at %s has been confirmed by the salon.", serviceName, salonName, apptDateStr, startTimeStr),
			Icon:  "/vite.svg",
			URL:   "/appointments",
		})
	}

	// 3. Schedule look-ahead reminders (if within next 70 mins)
	if h.BookingService != nil && h.BookingService.Scheduler != nil {
		st, _ := time.Parse("15:04", startTimeStr)
		dt, _ := time.Parse("2006-01-02", apptDateStr)
		apptTime := time.Date(dt.Year(), dt.Month(), dt.Day(), st.Hour(), st.Minute(), 0, 0, time.Local)

		if time.Until(apptTime) < 70*time.Minute {
			// Schedule 1h reminder
			h.BookingService.Scheduler.ScheduleReminder(context.Background(), apptTime, 60, appointmentID, customerID, serviceName, salonName, "reminder_1h_sent", "Appointment in 1 Hour", "Your appointment for %s at %s is in 1 hour. Get ready! 💇")

			// Schedule 20m reminder
			h.BookingService.Scheduler.ScheduleReminder(context.Background(), apptTime, 20, appointmentID, customerID, serviceName, salonName, "reminder_20m_sent", "Appointment in 20 Minutes", "Your appointment for %s at %s starts in 20 minutes. Head over! 🏃")
		}
	}
}

func (h *AppointmentHandler) CompleteAppointment(c *gin.Context) {
	appointmentID := c.Param("id")

	_, err := h.DB.Exec(context.Background(),
		"UPDATE appointments SET status = 'completed', updated_at = NOW() WHERE id = $1",
		appointmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete appointment"})
		return
	}

	// Mark payment as completed
	h.DB.Exec(context.Background(),
		"UPDATE payments SET status = 'completed', updated_at = NOW() WHERE appointment_id = $1",
		appointmentID)

	c.JSON(http.StatusOK, gin.H{"message": "Appointment completed"})
}

func (h *AppointmentHandler) AdjustAppointmentTime(c *gin.Context) {
	appointmentID := c.Param("id")
	var req models.AdjustAppointmentTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.BookingService.CascadeReschedule(c.Request.Context(), appointmentID, req.NewEndTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment adjusted and subsequent ones rescheduled"})
}

func (h *AppointmentHandler) GetAvailableSlots(c *gin.Context) {
	staffID := c.Query("staff_id")
	serviceID := c.Query("service_id")
	date := c.Query("date")

	if staffID == "" || serviceID == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "staff_id, service_id, and date query params required"})
		return
	}

	slots, err := h.BookingService.GetAvailableSlots(c.Request.Context(), staffID, serviceID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, slots)
}
