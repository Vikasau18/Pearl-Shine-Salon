package services

import (
	"context"
	"fmt"
	"time"

	"saloon-backend/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingService struct {
	DB *pgxpool.Pool
}

func NewBookingService(db *pgxpool.Pool) *BookingService {
	return &BookingService{DB: db}
}

// BookAppointment handles concurrency-safe booking using database-level locking
func (s *BookingService) BookAppointment(ctx context.Context, customerID string, req models.BookAppointmentRequest) (*models.Appointment, *models.Payment, error) {
	// Get service details for computing end time
	var durationMinutes, bufferMinutes int
	var servicePrice float64
	var serviceName string
	err := s.DB.QueryRow(ctx,
		"SELECT name, duration_minutes, buffer_minutes, price FROM services WHERE id = $1 AND is_active = true",
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
		 VALUES ($1, $2, $3, $4, $5, $6, $7, 'confirmed', $8, $9)
		 RETURNING id, customer_id, salon_id, staff_id, service_id, appointment_date::text, start_time::text, end_time::text, status, COALESCE(notes,''), promo_code_id, created_at, updated_at`,
		customerID, req.SalonID, req.StaffID, req.ServiceID, req.Date, req.StartTime, endTimeStr, req.Notes, promoCodeID,
	).Scan(&appt.ID, &appt.CustomerID, &appt.SalonID, &appt.StaffID, &appt.ServiceID,
		&appt.AppointmentDate, &appt.StartTime, &appt.EndTime, &appt.Status,
		&appt.Notes, &appt.PromoCodeID, &appt.CreatedAt, &appt.UpdatedAt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	// Create payment record
	tax := (servicePrice - discount) * 0.08 // 8% tax
	total := servicePrice - discount + tax
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

	appt.ServiceName = serviceName
	appt.ServicePrice = servicePrice

	return &appt, &payment, nil
}

// GetAvailableSlots returns available time slots for a staff member on a given date
func (s *BookingService) GetAvailableSlots(ctx context.Context, staffID, serviceID, dateStr string) ([]models.TimeSlot, error) {
	// Get service duration
	var durationMinutes int
	err := s.DB.QueryRow(ctx,
		"SELECT duration_minutes FROM services WHERE id = $1", serviceID).Scan(&durationMinutes)
	if err != nil {
		return nil, fmt.Errorf("service not found")
	}

	// Get working hours for this day
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date")
	}
	dayOfWeek := int(date.Weekday())

	// Debug log
	fmt.Printf("Checking slots for Staff: %s, Date: %s, DayOfWeek: %d\n", staffID, dateStr, dayOfWeek)

	var isOff bool
	var workStartStr, workEndStr string
	err = s.DB.QueryRow(ctx,
		`SELECT is_off, start_time::text, end_time::text 
		 FROM staff_working_hours WHERE staff_id = $1 AND day_of_week = $2`,
		staffID, dayOfWeek).Scan(&isOff, &workStartStr, &workEndStr)

	if err != nil {
		fmt.Printf("Error requesting working hours: %v\n", err)
		return []models.TimeSlot{}, nil
	}
	if isOff {
		fmt.Printf("Staff is off on day %d\n", dayOfWeek)
		return []models.TimeSlot{}, nil
	}

	fmt.Printf("Work hours: %s - %s\n", workStartStr, workEndStr)

	parseTime := func(t string) (time.Time, error) {
		if len(t) == 5 {
			return time.Parse("15:04", t)
		}
		return time.Parse("15:04:05", t)
	}

	workStart, err1 := parseTime(workStartStr)
	workEnd, err2 := parseTime(workEndStr)

	if err1 != nil || err2 != nil {
		fmt.Printf("Error parsing time: %v, %v\n", err1, err2)
		return []models.TimeSlot{}, fmt.Errorf("invalid working hours format")
	}

	// Get existing appointments for this staff on this date
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

	// Generate 30-minute interval slots
	var slots []models.TimeSlot
	slotDuration := time.Duration(durationMinutes) * time.Minute

	for t := workStart; t.Add(slotDuration).Before(workEnd) || t.Add(slotDuration).Equal(workEnd); t = t.Add(30 * time.Minute) {
		slotEnd := t.Add(slotDuration)
		available := true

		for _, b := range booked {
			// Check overlap
			if t.Before(b.end) && slotEnd.After(b.start) {
				available = false
				break
			}
		}

		slots = append(slots, models.TimeSlot{
			StartTime: t.Format("15:04"),
			EndTime:   slotEnd.Format("15:04"),
			Available: available,
		})
	}

	if slots == nil {
		slots = []models.TimeSlot{}
	}
	return slots, nil
}
