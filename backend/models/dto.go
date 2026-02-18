package models

// Request / Response DTOs

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateSalonRequest struct {
	Name        string  `json:"name" binding:"required"`
	Address     string  `json:"address" binding:"required"`
	City        string  `json:"city" binding:"required"`
	State       string  `json:"state"`
	ZipCode     string  `json:"zip_code"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	OpeningTime string  `json:"opening_time"`
	ClosingTime string  `json:"closing_time"`
}

type CreateServiceRequest struct {
	Name            string  `json:"name" binding:"required"`
	Description     string  `json:"description"`
	Price           float64 `json:"price" binding:"required"`
	DurationMinutes int     `json:"duration_minutes" binding:"required"`
	BufferMinutes   int     `json:"buffer_minutes"`
	Category        string  `json:"category"`
}

type CreateStaffRequest struct {
	Name           string   `json:"name" binding:"required"`
	Role           string   `json:"role"`
	Specialization string   `json:"specialization"`
	Bio            string   `json:"bio"`
	AvatarURL      string   `json:"avatar_url"`
	ServiceIDs     []string `json:"service_ids"`
}

type UpdateStaffWorkingHoursRequest struct {
	Hours []StaffWorkingHours `json:"hours" binding:"required"`
}

type BookAppointmentRequest struct {
	SalonID   string `json:"salon_id" binding:"required"`
	StaffID   string `json:"staff_id" binding:"required"`
	ServiceID string `json:"service_id" binding:"required"`
	Date      string `json:"date" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	Notes     string `json:"notes"`
	PromoCode string `json:"promo_code"`
}

type RescheduleRequest struct {
	Date      string `json:"date" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
}

type CreateReviewRequest struct {
	SalonID       string `json:"salon_id" binding:"required"`
	StaffID       string `json:"staff_id"`
	AppointmentID string `json:"appointment_id"`
	Rating        int    `json:"rating" binding:"required,min=1,max=5"`
	Comment       string `json:"comment"`
}

type ProcessPaymentRequest struct {
	AppointmentID string  `json:"appointment_id" binding:"required"`
	Method        string  `json:"method" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
}

type CreatePromoRequest struct {
	Code            string  `json:"code" binding:"required"`
	DiscountPercent float64 `json:"discount_percent" binding:"required"`
	ValidFrom       string  `json:"valid_from"`
	ValidUntil      string  `json:"valid_until"`
	MaxUses         *int    `json:"max_uses"`
}

type JoinWaitlistRequest struct {
	SalonID       string `json:"salon_id" binding:"required"`
	ServiceID     string `json:"service_id" binding:"required"`
	StaffID       string `json:"staff_id"`
	PreferredDate string `json:"preferred_date" binding:"required"`
	PreferredTime string `json:"preferred_time"`
}

type AvailableSlotsRequest struct {
	SalonID   string `form:"salon_id" binding:"required"`
	StaffID   string `form:"staff_id" binding:"required"`
	ServiceID string `form:"service_id" binding:"required"`
	Date      string `form:"date" binding:"required"`
}

type TimeSlot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Available bool   `json:"available"`
}

type AnalyticsResponse struct {
	TotalAppointments    int              `json:"total_appointments"`
	TotalRevenue         float64          `json:"total_revenue"`
	AverageRating        float64          `json:"average_rating"`
	PopularServices      []ServiceStats   `json:"popular_services"`
	PeakHours            []HourStats      `json:"peak_hours"`
	RevenueByMonth       []MonthlyRevenue `json:"revenue_by_month"`
	AppointmentsByStatus map[string]int   `json:"appointments_by_status"`
}

type ServiceStats struct {
	ServiceName string  `json:"service_name"`
	BookCount   int     `json:"book_count"`
	Revenue     float64 `json:"revenue"`
}

type HourStats struct {
	Hour      int `json:"hour"`
	BookCount int `json:"book_count"`
}

type MonthlyRevenue struct {
	Month   string  `json:"month"`
	Revenue float64 `json:"revenue"`
}

type PaginationParams struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalCount int         `json:"total_count"`
	TotalPages int         `json:"total_pages"`
}
