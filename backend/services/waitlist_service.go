package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WaitlistService struct {
	DB          *pgxpool.Pool
	PushService *PushService
}

func NewWaitlistService(db *pgxpool.Pool, ps *PushService) *WaitlistService {
	return &WaitlistService{DB: db, PushService: ps}
}

func (s *WaitlistService) AutoAssignWaitlist(ctx context.Context, salonID, serviceID, staffID, dateStr, startTime string) {
	if s.PushService == nil {
		return
	}

	fmt.Printf("Waitlist: Attempting auto-assign for salon %s, date %s, time %s\n", salonID, dateStr, startTime)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		fmt.Printf("Waitlist: Failed to start transaction: %v\n", err)
		return
	}
	defer tx.Rollback(ctx)

	// 1. Find the earliest customer waiting for this exact slot
	var waitlistID, customerID, reqServiceID string
	err = tx.QueryRow(ctx,
		`SELECT id, customer_id, service_id 
		 FROM waitlist 
		 WHERE salon_id = $1 AND preferred_date = $2 AND preferred_time = $3 AND status = 'waiting'
		 ORDER BY created_at ASC LIMIT 1`,
		salonID, dateStr, startTime).Scan(&waitlistID, &customerID, &reqServiceID)

	if err != nil {
		fmt.Printf("Waitlist: No matching waiting customer found for this slot.\n")
		return
	}

	// 2. Fetch service details for appointment creation
	var serviceName string
	var servicePrice float64
	var duration int
	err = tx.QueryRow(ctx, "SELECT name, price, duration_minutes FROM services WHERE id = $1", reqServiceID).
		Scan(&serviceName, &servicePrice, &duration)
	if err != nil {
		fmt.Printf("Waitlist: Failed to fetch service details: %v\n", err)
		return
	}

	// Calculate end time
	parsedStart, _ := time.Parse("15:04:05", startTime)
	if parsedStart.IsZero() {
		// Try without seconds if stored differently
		parsedStart, _ = time.Parse("15:04", startTime)
	}
	endTime := parsedStart.Add(time.Duration(duration) * time.Minute).Format("15:04:05")

	// 3. Create the appointment
	var apptID string
	err = tx.QueryRow(ctx,
		`INSERT INTO appointments (customer_id, salon_id, staff_id, service_id, appointment_date, start_time, end_time, status, notes)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, 'pending', 'Auto-assigned from waitlist')
		 RETURNING id`,
		customerID, salonID, staffID, reqServiceID, dateStr, startTime, endTime).Scan(&apptID)
	if err != nil {
		fmt.Printf("Waitlist: Failed to create auto-appointment: %v\n", err)
		return
	}

	// 4. Create payment record (0 tax as per previous requirement)
	receiptNumber := fmt.Sprintf("RCP-AUTO-%s-%s", time.Now().Format("20060102"), apptID[:8])
	_, err = tx.Exec(ctx,
		`INSERT INTO payments (appointment_id, amount, discount, tax, total, method, status, receipt_number)
		 VALUES ($1, $2, 0, 0, $2, 'cash', 'pending', $3)`,
		apptID, servicePrice, receiptNumber)
	if err != nil {
		fmt.Printf("Waitlist: Failed to create payment: %v\n", err)
		return
	}

	// 5. Mark all waitlist entries for this user on this date as 'booked'
	_, err = tx.Exec(ctx,
		"UPDATE waitlist SET status = 'booked' WHERE customer_id = $1 AND preferred_date = $2 AND salon_id = $3",
		customerID, dateStr, salonID)
	if err != nil {
		fmt.Printf("Waitlist: Failed to update waitlist status: %v\n", err)
		return
	}

	// 6. Check for an existing "safe" appointment to cancel (Slot Swap)
	var oldApptID, oldStartTime string
	err = tx.QueryRow(ctx,
		`SELECT id, start_time::text FROM appointments 
		 WHERE customer_id = $1 AND appointment_date = $2 AND salon_id = $3 
		 AND status IN ('pending', 'confirmed') AND id != $4 LIMIT 1`,
		customerID, dateStr, salonID, apptID).Scan(&oldApptID, &oldStartTime)

	swapped := false
	if err == nil {
		_, err = tx.Exec(ctx,
			"UPDATE appointments SET status = 'cancelled', notes = $1 WHERE id = $2",
			fmt.Sprintf("Cancelled in favor of preferred slot at %s", startTime[:5]), oldApptID)
		if err == nil {
			swapped = true
		}
	}

	if err := tx.Commit(ctx); err != nil {
		fmt.Printf("Waitlist: Failed to commit auto-assignment: %v\n", err)
		return
	}

	// 7. Notify the user
	var salonName string
	s.DB.QueryRow(ctx, "SELECT name FROM salons WHERE id = $1", salonID).Scan(&salonName)

	body := fmt.Sprintf("You've been auto-assigned to the %s slot at %s on %s!", startTime[:5], salonName, dateStr)
	if swapped {
		body = fmt.Sprintf("Moved to your preferred %s slot at %s! Your original %s booking was cancelled.", startTime[:5], salonName, oldStartTime[:5])
	}

	go s.PushService.SendToUser(context.Background(), customerID, PushPayload{
		Title: "Slot Assigned 💇",
		Body:  body,
		Icon:  "/vite.svg",
		URL:   "/appointments",
	})

	fmt.Printf("Waitlist: Successfully auto-assigned user %s to slot %s (swapped: %v)\n", customerID, startTime, swapped)
}
