package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Scheduler struct {
	DB   *pgxpool.Pool
	Push *PushService
}

func NewScheduler(db *pgxpool.Pool, push *PushService) *Scheduler {
	return &Scheduler{DB: db, Push: push}
}

// Start begins the background reminder loop. Call this in a goroutine.
func (s *Scheduler) Start(ctx context.Context) {
	log.Println("📅 Appointment reminder scheduler started (1-hour optimized mode)")
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// Run immediately on start
	s.checkReminders(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("📅 Scheduler stopped")
			return
		case <-ticker.C:
			s.checkReminders(ctx)
		}
	}
}

func (s *Scheduler) checkReminders(ctx context.Context) {
	// Look ahead for appointments in the next 1 hour + 5 min buffer
	s.queryAndSchedule(ctx, 65)
}

func (s *Scheduler) queryAndSchedule(ctx context.Context, lookaheadMinutes int) {
	// Query for appointments that need EITHER 1h or 20m reminders in the next window
	query := `
		SELECT a.id, a.customer_id, a.appointment_date, a.start_time, 
		       s.name AS service_name, sal.name AS salon_name,
		       a.reminder_1h_sent, a.reminder_20m_sent
		FROM appointments a
		JOIN services s ON s.id = a.service_id
		JOIN salons sal ON sal.id = a.salon_id
		WHERE a.status IN ('confirmed', 'pending')
		  AND (a.reminder_1h_sent = false OR a.reminder_20m_sent = false)
		  AND (a.appointment_date + a.start_time) BETWEEN 
		      NOW() AND (NOW() + INTERVAL '1 minute' * $1)
	`
	rows, err := s.DB.Query(ctx, query, lookaheadMinutes)
	if err != nil {
		log.Printf("Scheduler: query error: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var apptID, customerID, serviceName, salonName string
		var apptDate time.Time
		var startTimeStr string
		var sent1h, sent20m bool

		if err := rows.Scan(&apptID, &customerID, &apptDate, &startTimeStr, &serviceName, &salonName, &sent1h, &sent20m); err != nil {
			log.Printf("Scheduler: scan error: %v", err)
			continue
		}

		// Parse full appointment time
		format := "15:04:05"
		if len(startTimeStr) == 5 {
			format = "15:04"
		}
		st, _ := time.Parse(format, startTimeStr)
		apptTime := time.Date(apptDate.Year(), apptDate.Month(), apptDate.Day(), st.Hour(), st.Minute(), 0, 0, time.Local)

		// Schedule 1h reminder
		if !sent1h {
			s.ScheduleReminder(ctx, apptTime, 60, apptID, customerID, serviceName, salonName, "reminder_1h_sent", "Appointment in 1 Hour", "Your appointment for %s at %s is in 1 hour. Get ready! 💇")
		}

		// Schedule 20m reminder
		if !sent20m {
			s.ScheduleReminder(ctx, apptTime, 20, apptID, customerID, serviceName, salonName, "reminder_20m_sent", "Appointment in 20 Minutes", "Your appointment for %s at %s starts in 20 minutes. Head over! 🏃")
		}
	}
}

func (s *Scheduler) ScheduleReminder(ctx context.Context, apptTime time.Time, minBefore int, apptID, customerID, serviceName, salonName, column, title, template string) {
	reminderTime := apptTime.Add(-time.Duration(minBefore) * time.Minute)
	delay := time.Until(reminderTime)

	// If the time has already passed but is very recent (e.g. within 1 min), send it now
	if delay <= 0 && delay > -1*time.Minute {
		delay = 0
	}

	if delay >= 0 {
		time.AfterFunc(delay, func() {
			s.executeReminder(context.Background(), apptID, customerID, serviceName, salonName, column, title, template)
		})
	}
}

func (s *Scheduler) executeReminder(ctx context.Context, apptID, customerID, serviceName, salonName, column, title, template string) {
	// Re-verify the appointment hasn't been cancelled or already sent by another process
	var alreadySent bool
	err := s.DB.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM appointments WHERE id = $1 AND status NOT IN ('cancelled', 'no_show')", column), apptID).Scan(&alreadySent)
	if err != nil || alreadySent {
		return
	}

	body := fmt.Sprintf(template, serviceName, salonName)

	// 1. Send Push
	s.Push.SendToUser(ctx, customerID, PushPayload{
		Title: title,
		Body:  body,
		Icon:  "/vite.svg",
		URL:   "/appointments",
	})

	// 2. Create In-App
	s.DB.Exec(ctx,
		`INSERT INTO notifications (user_id, type, title, message, appointment_id)
		 VALUES ($1, 'appointment_reminder', $2, $3, $4)`,
		customerID, title, body, apptID)

	// 3. Mark Sent
	s.DB.Exec(ctx, fmt.Sprintf("UPDATE appointments SET %s = true WHERE id = $1", column), apptID)

	log.Printf("📅 Reminder sent (%s): %s for user %s", column, serviceName, customerID[:8])
}
