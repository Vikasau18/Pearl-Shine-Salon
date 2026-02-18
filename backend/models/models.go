package models

import (
	"time"
)

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"`
	Name          string    `json:"name"`
	Phone         string    `json:"phone,omitempty"`
	Role          string    `json:"role"`
	AvatarURL     string    `json:"avatar_url,omitempty"`
	GoogleID      string    `json:"-"`
	LoyaltyPoints int       `json:"loyalty_points"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Salon struct {
	ID           string    `json:"id"`
	OwnerID      string    `json:"owner_id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	City         string    `json:"city"`
	State        string    `json:"state,omitempty"`
	ZipCode      string    `json:"zip_code,omitempty"`
	Lat          float64   `json:"lat,omitempty"`
	Lng          float64   `json:"lng,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	Email        string    `json:"email,omitempty"`
	Description  string    `json:"description,omitempty"`
	Rating       float64   `json:"rating"`
	TotalReviews int       `json:"total_reviews"`
	ImageURL     string    `json:"image_url,omitempty"`
	IsActive     bool      `json:"is_active"`
	OpeningTime  string    `json:"opening_time"`
	ClosingTime  string    `json:"closing_time"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	// Computed / Joined fields
	Distance    *float64 `json:"distance,omitempty"`
	IsFavorited *bool    `json:"is_favorited,omitempty"`
}

type SalonGallery struct {
	ID        string    `json:"id"`
	SalonID   string    `json:"salon_id"`
	ImageURL  string    `json:"image_url"`
	Caption   string    `json:"caption,omitempty"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

type Staff struct {
	ID             string    `json:"id"`
	SalonID        string    `json:"salon_id"`
	UserID         *string   `json:"user_id,omitempty"`
	Name           string    `json:"name"`
	Role           string    `json:"role"`
	Specialization string    `json:"specialization,omitempty"`
	Bio            string    `json:"bio,omitempty"`
	AvatarURL      string    `json:"avatar_url,omitempty"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	// Joined
	Services []Service `json:"services,omitempty"`
}

type StaffWorkingHours struct {
	ID        string `json:"id"`
	StaffID   string `json:"staff_id"`
	DayOfWeek int    `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	IsOff     bool   `json:"is_off"`
}

type Service struct {
	ID              string    `json:"id"`
	SalonID         string    `json:"salon_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	Price           float64   `json:"price"`
	DurationMinutes int       `json:"duration_minutes"`
	BufferMinutes   int       `json:"buffer_minutes"`
	Category        string    `json:"category,omitempty"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Appointment struct {
	ID              string    `json:"id"`
	CustomerID      string    `json:"customer_id"`
	SalonID         string    `json:"salon_id"`
	StaffID         string    `json:"staff_id"`
	ServiceID       string    `json:"service_id"`
	AppointmentDate string    `json:"appointment_date"`
	StartTime       string    `json:"start_time"`
	EndTime         string    `json:"end_time"`
	Status          string    `json:"status"`
	Notes           string    `json:"notes,omitempty"`
	PromoCodeID     *string   `json:"promo_code_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	// Joined
	SalonName     string  `json:"salon_name,omitempty"`
	StaffName     string  `json:"staff_name,omitempty"`
	ServiceName   string  `json:"service_name,omitempty"`
	ServicePrice  float64 `json:"service_price,omitempty"`
	TotalPrice    float64 `json:"total_price,omitempty"`
	CustomerName  string  `json:"customer_name,omitempty"`
	CustomerEmail string  `json:"customer_email,omitempty"`
}

type Payment struct {
	ID            string    `json:"id"`
	AppointmentID string    `json:"appointment_id"`
	Amount        float64   `json:"amount"`
	Discount      float64   `json:"discount"`
	Tax           float64   `json:"tax"`
	Total         float64   `json:"total"`
	Method        string    `json:"method"`
	Status        string    `json:"status"`
	ReceiptNumber string    `json:"receipt_number"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Review struct {
	ID            string    `json:"id"`
	CustomerID    string    `json:"customer_id"`
	SalonID       string    `json:"salon_id"`
	StaffID       *string   `json:"staff_id,omitempty"`
	AppointmentID *string   `json:"appointment_id,omitempty"`
	Rating        int       `json:"rating"`
	Comment       string    `json:"comment,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	// Joined
	CustomerName   string `json:"customer_name,omitempty"`
	CustomerAvatar string `json:"customer_avatar,omitempty"`
}

type Favorite struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	SalonID   string    `json:"salon_id"`
	CreatedAt time.Time `json:"created_at"`
	// Joined
	Salon *Salon `json:"salon,omitempty"`
}

type PromoCode struct {
	ID              string     `json:"id"`
	SalonID         string     `json:"salon_id"`
	Code            string     `json:"code"`
	DiscountPercent float64    `json:"discount_percent"`
	ValidFrom       time.Time  `json:"valid_from"`
	ValidUntil      *time.Time `json:"valid_until,omitempty"`
	MaxUses         *int       `json:"max_uses,omitempty"`
	UsedCount       int        `json:"used_count"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
}

type Waitlist struct {
	ID            string    `json:"id"`
	CustomerID    string    `json:"customer_id"`
	SalonID       string    `json:"salon_id"`
	ServiceID     string    `json:"service_id"`
	StaffID       *string   `json:"staff_id,omitempty"`
	PreferredDate string    `json:"preferred_date"`
	PreferredTime *string   `json:"preferred_time,omitempty"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type Notification struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Type          string    `json:"type"`
	Title         string    `json:"title"`
	Message       string    `json:"message"`
	IsRead        bool      `json:"is_read"`
	AppointmentID *string   `json:"appointment_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
