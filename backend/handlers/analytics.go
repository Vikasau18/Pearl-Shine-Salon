package handlers

import (
	"context"
	"net/http"

	"saloon-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticsHandler struct {
	DB *pgxpool.Pool
}

func NewAnalyticsHandler(db *pgxpool.Pool) *AnalyticsHandler {
	return &AnalyticsHandler{DB: db}
}

func (h *AnalyticsHandler) GetDashboard(c *gin.Context) {
	salonID := c.Param("salon_id")

	var analytics models.AnalyticsResponse

	// Total appointments
	h.DB.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM appointments WHERE salon_id = $1", salonID).Scan(&analytics.TotalAppointments)

	// Total revenue
	h.DB.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(p.total), 0) FROM payments p 
		 JOIN appointments a ON a.id = p.appointment_id 
		 WHERE a.salon_id = $1 AND p.status = 'completed'`, salonID).Scan(&analytics.TotalRevenue)

	// Average rating
	h.DB.QueryRow(context.Background(),
		"SELECT COALESCE(AVG(rating), 0) FROM reviews WHERE salon_id = $1", salonID).Scan(&analytics.AverageRating)

	// Appointments by status
	analytics.AppointmentsByStatus = make(map[string]int)
	statusRows, _ := h.DB.Query(context.Background(),
		"SELECT status, COUNT(*) FROM appointments WHERE salon_id = $1 GROUP BY status", salonID)
	if statusRows != nil {
		for statusRows.Next() {
			var status string
			var count int
			statusRows.Scan(&status, &count)
			analytics.AppointmentsByStatus[status] = count
		}
		statusRows.Close()
	}

	// Popular services
	svcRows, _ := h.DB.Query(context.Background(),
		`SELECT sv.name, COUNT(*) as book_count, COALESCE(SUM(p.total), 0) as revenue
		 FROM appointments a
		 JOIN services sv ON sv.id = a.service_id
		 LEFT JOIN payments p ON p.appointment_id = a.id AND p.status = 'completed'
		 WHERE a.salon_id = $1
		 GROUP BY sv.name
		 ORDER BY book_count DESC
		 LIMIT 10`, salonID)
	if svcRows != nil {
		for svcRows.Next() {
			var s models.ServiceStats
			svcRows.Scan(&s.ServiceName, &s.BookCount, &s.Revenue)
			analytics.PopularServices = append(analytics.PopularServices, s)
		}
		svcRows.Close()
	}
	if analytics.PopularServices == nil {
		analytics.PopularServices = []models.ServiceStats{}
	}

	// Peak hours
	hourRows, _ := h.DB.Query(context.Background(),
		`SELECT EXTRACT(HOUR FROM start_time)::int as hour, COUNT(*) as book_count
		 FROM appointments WHERE salon_id = $1 AND status NOT IN ('cancelled', 'no_show')
		 GROUP BY hour ORDER BY hour`, salonID)
	if hourRows != nil {
		for hourRows.Next() {
			var h models.HourStats
			hourRows.Scan(&h.Hour, &h.BookCount)
			analytics.PeakHours = append(analytics.PeakHours, h)
		}
		hourRows.Close()
	}
	if analytics.PeakHours == nil {
		analytics.PeakHours = []models.HourStats{}
	}

	// Revenue by month
	monthRows, _ := h.DB.Query(context.Background(),
		`SELECT TO_CHAR(p.created_at, 'YYYY-MM') as month, SUM(p.total) as revenue
		 FROM payments p
		 JOIN appointments a ON a.id = p.appointment_id
		 WHERE a.salon_id = $1 AND p.status = 'completed'
		 GROUP BY month ORDER BY month DESC LIMIT 12`, salonID)
	if monthRows != nil {
		for monthRows.Next() {
			var m models.MonthlyRevenue
			monthRows.Scan(&m.Month, &m.Revenue)
			analytics.RevenueByMonth = append(analytics.RevenueByMonth, m)
		}
		monthRows.Close()
	}
	if analytics.RevenueByMonth == nil {
		analytics.RevenueByMonth = []models.MonthlyRevenue{}
	}

	c.JSON(http.StatusOK, analytics)
}
