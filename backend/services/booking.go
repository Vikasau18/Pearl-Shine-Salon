package services

import (
	"context"
	"fmt"
	"time"

	"saloon-backend/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingService struct {
	DB          *pgxpool.Pool
	Scheduler   *Scheduler
	PushService *PushService
}

func NewBookingService(db *pgxpool.Pool, scheduler *Scheduler) *BookingService {
	return &BookingService{DB: db, Scheduler: scheduler}
}

func (s *BookingService) SetPushService(ps *PushService) {
	s.PushService = ps
}

// CascadeReschedule adjusts the target appointment and ripples the shift to subsequent ones
func (s *BookingService) CascadeReschedule(ctx context.Context, apptID string, newEndTimeStr string) error {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Fetch the targeted appointment
	var staffID, dateStr, currentStartTime string
	err = tx.QueryRow(ctx,
		"SELECT staff_id, appointment_date::text, start_time::text FROM appointments WHERE id = $1",
		apptID).Scan(&staffID, &dateStr, &currentStartTime)
	if err != nil {
		return fmt.Errorf("appointment not found")
	}

	// Update the target appointment's end time
	_, err = tx.Exec(ctx, "UPDATE appointments SET end_time = $1 WHERE id = $2", newEndTimeStr, apptID)
	if err != nil {
		return fmt.Errorf("failed to update target end time")
	}

	// 2. Fetch all subsequent appointments for this staff member today
	rows, err := tx.Query(ctx,
		`SELECT id, customer_id, start_time::text, end_time::text, service_id 
		 FROM appointments 
		 WHERE staff_id = $1 AND appointment_date = $2 
		 AND status IN ('pending', 'confirmed') 
		 AND start_time > $3::time 
		 ORDER BY start_time ASC`,
		staffID, dateStr, currentStartTime)
	if err != nil {
		return err
	}
	defer rows.Close()

	type apptShift struct {
		id, customerID, serviceID string
		oldStart, oldEnd          string
		duration                  time.Duration
	}
	var subsequent []apptShift

	for rows.Next() {
		var a apptShift
		var startStr, endStr string
		rows.Scan(&a.id, &a.customerID, &startStr, &endStr, &a.serviceID)

		st, _ := time.Parse("15:04:05", startStr)
		en, _ := time.Parse("15:04:05", endStr)
		a.duration = en.Sub(st)
		a.oldStart = startStr
		subsequent = append(subsequent, a)
	}
	rows.Close()

	// 3. Ripple the changes
	lastEndTime, _ := time.Parse("15:04", newEndTimeStr)
	if lastEndTime.IsZero() {
		lastEndTime, _ = time.Parse("15:04:05", newEndTimeStr)
	}

	buffer := 5 * time.Minute

	for _, a := range subsequent {
		newStart := lastEndTime.Add(buffer)
		newEnd := newStart.Add(a.duration)

		newStartStr := newStart.Format("15:04:05")
		newEndStr := newEnd.Format("15:04:05")

		_, err = tx.Exec(ctx,
			"UPDATE appointments SET start_time = $1, end_time = $2, notes = COALESCE(notes, '') || '\n(Time shifted due to preceding appointment delay)' WHERE id = $3",
			newStartStr, newEndStr, a.id)
		if err != nil {
			return fmt.Errorf("failed to update subsequent appointment %s", a.id)
		}

		// Notify the user immediately if PushService is available
		if s.PushService != nil {
			go s.PushService.SendToUser(context.Background(), a.customerID, PushPayload{
				Title: "Appointment Update 🕒",
				Body:  fmt.Sprintf("Your appointment today has been shifted to %s due to a slight delay. We apologize for the inconvenience!", newStart.Format("15:04")),
				Icon:  "/vite.svg",
				URL:   "/appointments",
			})
		}

		lastEndTime = newEnd
	}

	return tx.Commit(ctx)
}

// BookAppointment handles concurrency-safe booking using database-level locking
func (s *BookingService) BookAppointment(ctx context.Context, customerID string, req models.BookAppointmentRequest) (*models.Appointment, *models.Payment, error) {
	// Get service details for computing end time
	var durationMinutes, bufferMinutes int
	var servicePrice float64
	var serviceName string
	err := s.DB.QueryRow(ctx,
		`SELECT s.name, s.duration_minutes, s.buffer_minutes, s.price 
		 FROM services s 
		 WHERE s.id = $1 AND s.is_active = true`,
		req.ServiceID).Scan(&serviceName, &durationMinutes, &bufferMinutes, &servicePrice)
	if err != nil {
		return nil, nil, fmt.Errorf("service not found")
	}

	// Parse start time and calculate end time
	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid start time format, use HH:MM")
	}
	endTime := startTime.Add(time.Duration(durationMinutes) * time.Minute)
	endTimeStr := endTime.Format("15:04")

	// Get salon name for notifications
	var salonName string
	err = s.DB.QueryRow(ctx, "SELECT name FROM salons WHERE id = $1", req.SalonID).Scan(&salonName)
	if err != nil {
		return nil, nil, fmt.Errorf("salon not found")
	}

	// Check staff working hours
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}
	dayOfWeek := int(date.Weekday())

	var isOff bool
	var workStart, workEnd string
	err = s.DB.QueryRow(ctx,
		`SELECT is_off, start_time::text, end_time::text 
		 FROM staff_working_hours WHERE staff_id = $1 AND day_of_week = $2`,
		req.StaffID, dayOfWeek).Scan(&isOff, &workStart, &workEnd)
	if err == nil && isOff {
		return nil, nil, fmt.Errorf("staff member is not available on this day")
	}

	// Begin transaction with row-level lock to prevent double booking
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check for salon closure
	var closureCount int
	err = tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM salon_closures
		 WHERE salon_id = $1
		 AND $2::date >= start_date AND $2::date <= end_date
		 AND (
			 start_time IS NULL OR 
			 end_time IS NULL OR 
			 (start_time::time < $4::time AND end_time::time > $3::time)
		 )`,
		req.SalonID, req.Date, req.StartTime, endTimeStr).Scan(&closureCount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check closures: %w", err)
	}
	if closureCount > 0 {
		return nil, nil, fmt.Errorf("salon is closed during this time")
	}

	// Lock the staff row to serialize bookings for this staff member
	// This prevents race conditions where two concurrent requests see the slot as available
	_, err = tx.Exec(ctx, "SELECT 1 FROM staff WHERE id = $1 FOR UPDATE", req.StaffID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to lock staff row: %w", err)
	}

	// Lock overlapping appointment rows for this staff on this date
	// This prevents concurrent bookings from creating overlapping slots
	var conflictCount int
	err = tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM appointments 
		 WHERE staff_id = $1 AND appointment_date = $2 
		 AND status NOT IN ('cancelled', 'no_show')
		 AND (start_time, end_time) OVERLAPS ($3::time, $4::time)`,
		req.StaffID, req.Date, req.StartTime, endTimeStr).Scan(&conflictCount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check availability: %w", err)
	}
	if conflictCount > 0 {
		return nil, nil, fmt.Errorf("time slot is already booked")
	}

	// Handle promo code
	var promoCodeID *string
	var discount float64
	if req.PromoCode != "" {
		var pc models.PromoCode
		err = tx.QueryRow(ctx,
			`SELECT id, discount_percent, max_uses, used_count FROM promo_codes 
			 WHERE code = $1 AND salon_id = $2 AND is_active = true 
			 AND (valid_until IS NULL OR valid_until > NOW())
			 FOR UPDATE`,
			req.PromoCode, req.SalonID).Scan(&pc.ID, &pc.DiscountPercent, &pc.MaxUses, &pc.UsedCount)
		if err == nil {
			if pc.MaxUses != nil && pc.UsedCount >= *pc.MaxUses {
				return nil, nil, fmt.Errorf("promo code has reached maximum uses")
			}
			promoCodeID = &pc.ID
			discount = servicePrice * (pc.DiscountPercent / 100.0)

			tx.Exec(ctx, "UPDATE promo_codes SET used_count = used_count + 1 WHERE id = $1", pc.ID)
		}
	}

	// Insert appointment
	var appt models.Appointment
	err = tx.QueryRow(ctx,
		`INSERT INTO appointments (customer_id, salon_id, staff_id, service_id, appointment_date, start_time, end_time, status, notes, promo_code_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, 'pending', $8, $9)
		 RETURNING id, customer_id, salon_id, staff_id, service_id, appointment_date::text, start_time::text, end_time::text, status, COALESCE(notes,''), promo_code_id, created_at, updated_at`,
		customerID, req.SalonID, req.StaffID, req.ServiceID, req.Date, req.StartTime, endTimeStr, req.Notes, promoCodeID,
	).Scan(&appt.ID, &appt.CustomerID, &appt.SalonID, &appt.StaffID, &appt.ServiceID,
		&appt.AppointmentDate, &appt.StartTime, &appt.EndTime, &appt.Status,
		&appt.Notes, &appt.PromoCodeID, &appt.CreatedAt, &appt.UpdatedAt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	// Create payment record
	tax := 0.0
	total := servicePrice - discount
	receiptNumber := fmt.Sprintf("RCP-%s-%s", time.Now().Format("20060102"), appt.ID[:8])

	var payment models.Payment
	err = tx.QueryRow(ctx,
		`INSERT INTO payments (appointment_id, amount, discount, tax, total, method, status, receipt_number) 
		 VALUES ($1, $2, $3, $4, $5, 'card', 'pending', $6)
		 RETURNING id, appointment_id, amount, discount, tax, total, method, status, receipt_number, created_at, updated_at`,
		appt.ID, servicePrice, discount, tax, total, receiptNumber,
	).Scan(&payment.ID, &payment.AppointmentID, &payment.Amount, &payment.Discount,
		&payment.Tax, &payment.Total, &payment.Method, &payment.Status,
		&payment.ReceiptNumber, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Create notification
	tx.Exec(ctx,
		`INSERT INTO notifications (user_id, type, title, message, appointment_id)
		 VALUES ($1, 'appointment_confirmed', 'Appointment Confirmed', $2, $3)`,
		customerID,
		fmt.Sprintf("Your appointment for %s on %s at %s has been confirmed.", serviceName, req.Date, req.StartTime),
		appt.ID)

	// Award loyalty points (1 point per dollar spent)
	tx.Exec(ctx,
		"UPDATE users SET loyalty_points = loyalty_points + $1 WHERE id = $2",
		int(servicePrice), customerID)

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to commit booking: %w", err)
	}

	// Parse full appointment time for scheduler
	st, _ := time.Parse("15:04", req.StartTime)
	dt, _ := time.Parse("2006-01-02", req.Date)
	apptTime := time.Date(dt.Year(), dt.Month(), dt.Day(), st.Hour(), st.Minute(), 0, 0, time.Local)

	// If the appointment is within the next 70 minutes, schedule reminders immediately
	// (Scheduler loop only runs hourly, so we handle mid-hour bookings here)
	if s.Scheduler != nil && time.Until(apptTime) < 70*time.Minute {
		// Schedule 1h reminder
		s.Scheduler.ScheduleReminder(ctx, apptTime, 60, appt.ID, customerID, serviceName, salonName, "reminder_1h_sent", "Appointment in 1 Hour", "Your appointment for %s at %s is in 1 hour. Get ready! 💇")

		// Schedule 20m reminder
		s.Scheduler.ScheduleReminder(ctx, apptTime, 20, appt.ID, customerID, serviceName, salonName, "reminder_20m_sent", "Appointment in 20 Minutes", "Your appointment for %s at %s starts in 20 minutes. Head over! 🏃")
	}

	return &appt, &payment, nil
}

// GetAvailableSlots returns time slots within SALON hours, marked with staff shift availability
func (s *BookingService) GetAvailableSlots(ctx context.Context, staffID, serviceID, dateStr string) ([]models.TimeSlot, error) {
	// 1. Get service duration and salon ID
	var durationMinutes int
	var salonID string
	err := s.DB.QueryRow(ctx,
		"SELECT duration_minutes, salon_id FROM services WHERE id = $1", serviceID).Scan(&durationMinutes, &salonID)
	if err != nil {
		return nil, fmt.Errorf("service not found")
	}

	// 2. Get SALON working hours (the overall grid boundaries)
	var salonStartStr, salonEndStr string
	err = s.DB.QueryRow(ctx,
		"SELECT opening_time::text, closing_time::text FROM salons WHERE id = $1", salonID).Scan(&salonStartStr, &salonEndStr)
	if err != nil {
		return nil, fmt.Errorf("salon not found")
	}

	// 2b. Check for salon closures on this date
	var closureStartTime, closureEndTime *string
	var isFullDayClosure bool
	closureRow := s.DB.QueryRow(ctx,
		`SELECT start_time::text, end_time::text
		 FROM salon_closures
		 WHERE salon_id = $1 AND $2::date >= start_date AND $2::date <= end_date
		 ORDER BY created_at DESC LIMIT 1`, salonID, dateStr)
	var cst, cet *string
	var cstVal, cetVal string
	if err := closureRow.Scan(&cstVal, &cetVal); err == nil {
		// Found a closure
		if cstVal == "" && cetVal == "" {
			isFullDayClosure = true
		} else {
			cst = &cstVal
			cet = &cetVal
			closureStartTime = cst
			closureEndTime = cet
		}
	} else {
		// Check for full-day closure with NULL times
		_ = s.DB.QueryRow(ctx,
			`SELECT 1 FROM salon_closures
			 WHERE salon_id = $1 AND $2::date >= start_date AND $2::date <= end_date
			 AND start_time IS NULL AND end_time IS NULL
			 LIMIT 1`, salonID, dateStr).Scan(new(int))
		// if scan succeeds, it's a full day closure; otherwise, no closure
		err2 := s.DB.QueryRow(ctx,
			`SELECT 1 FROM salon_closures
			 WHERE salon_id = $1 AND $2::date >= start_date AND $2::date <= end_date
			 AND start_time IS NULL
			 LIMIT 1`, salonID, dateStr).Scan(new(int))
		if err2 == nil {
			isFullDayClosure = true
		}
	}

	// If full-day closure, return a single "closed" marker
	if isFullDayClosure {
		return []models.TimeSlot{{
			StartTime:          salonStartStr,
			EndTime:            salonEndStr,
			Available:          false,
			AvailabilityReason: "closed",
		}}, nil
	}

	// 3. Get STAFF working hours for this specific day
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date")
	}
	dayOfWeek := int(date.Weekday())

	var staffOff bool
	var staffStartStr, staffEndStr string
	err = s.DB.QueryRow(ctx,
		`SELECT is_off, start_time::text, end_time::text 
		 FROM staff_working_hours WHERE staff_id = $1 AND day_of_week = $2`,
		staffID, dayOfWeek).Scan(&staffOff, &staffStartStr, &staffEndStr)

	// 4. Get existing appointments to check for "booked" slots
	rows, err := s.DB.Query(ctx,
		`SELECT start_time::text, end_time::text FROM appointments 
		 WHERE staff_id = $1 AND appointment_date = $2 
		 AND status NOT IN ('cancelled', 'no_show')
		 ORDER BY start_time`, staffID, dateStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type bookedSlot struct {
		start, end time.Time
	}
	var booked []bookedSlot
	for rows.Next() {
		var startStr, endStr string
		rows.Scan(&startStr, &endStr)
		st, _ := time.Parse("15:04:05", startStr)
		en, _ := time.Parse("15:04:05", endStr)
		booked = append(booked, bookedSlot{st, en})
	}

	// 5. Setup parsing and generation
	parseTime := func(t string) (time.Time, error) {
		if len(t) == 5 {
			return time.Parse("15:04", t)
		}
		return time.Parse("15:04:05", t)
	}

	salonStart, _ := parseTime(salonStartStr)
	salonEnd, _ := parseTime(salonEndStr)
	staffStart, _ := parseTime(staffStartStr)
	staffEnd, _ := parseTime(staffEndStr)

	// Parse partial closure times if present
	var closureStart, closureEnd time.Time
	hasPartialClosure := false
	if closureStartTime != nil && closureEndTime != nil {
		closureStart, _ = parseTime(*closureStartTime)
		closureEnd, _ = parseTime(*closureEndTime)
		hasPartialClosure = true
	}

	var slots []models.TimeSlot
	slotDuration := time.Duration(durationMinutes) * time.Minute

	// Today's filtering logic
	now := time.Now()
	isToday := date.Year() == now.Year() && date.Month() == now.Month() && date.Day() == now.Day()
	currentTimeStr := now.Format("15:04")
	currentTime, _ := time.Parse("15:04", currentTimeStr)

	// 6. Generate slots based on SALON hours
	for t := salonStart; t.Add(slotDuration).Before(salonEnd) || t.Add(slotDuration).Equal(salonEnd); t = t.Add(30 * time.Minute) {
		slotEnd := t.Add(slotDuration)

		// Filter past slots if today
		if isToday && t.Before(currentTime) {
			continue
		}

		available := true
		reason := "available"

		// 0. Check partial closure
		if hasPartialClosure && t.Before(closureEnd) && slotEnd.After(closureStart) {
			available = false
			reason = "closed"
		} else if !staffOff && (t.After(staffStart) || t.Equal(staffStart)) && (slotEnd.Before(staffEnd) || slotEnd.Equal(staffEnd)) {
			// A. Within staff shift — check for booked
			for _, b := range booked {
				if t.Before(b.end) && slotEnd.After(b.start) {
					available = false
					reason = "booked"
					break
				}
			}
		} else {
			available = false
			reason = "out_of_shift"
		}

		slots = append(slots, models.TimeSlot{
			StartTime:          t.Format("15:04"),
			EndTime:            slotEnd.Format("15:04"),
			Available:          available,
			AvailabilityReason: reason,
		})
	}

	if slots == nil {
		slots = []models.TimeSlot{}
	}
	return slots, nil
}
