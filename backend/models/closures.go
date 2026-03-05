package models

import "time"

type SalonClosure struct {
	ID        string    `json:"id"`
	SalonID   string    `json:"salon_id"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date"`
	StartTime *string   `json:"start_time,omitempty"` // nil = full day
	EndTime   *string   `json:"end_time,omitempty"`   // nil = full day
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateClosureRequest struct {
	StartDate string  `json:"start_date" binding:"required"`
	EndDate   string  `json:"end_date" binding:"required"`
	StartTime *string `json:"start_time"` // optional, for partial day
	EndTime   *string `json:"end_time"`   // optional, for partial day
	Reason    string  `json:"reason" binding:"required"`
}

type ConflictingAppointment struct {
	ID              string `json:"id"`
	CustomerName    string `json:"customer_name"`
	ServiceName     string `json:"service_name"`
	AppointmentDate string `json:"appointment_date"`
	StartTime       string `json:"start_time"`
	Status          string `json:"status"`
}

type CreateClosureResponse struct {
	Closure             SalonClosure             `json:"closure"`
	ConflictingCount    int                      `json:"conflicting_count"`
	ConflictingBookings []ConflictingAppointment `json:"conflicting_bookings"`
}
