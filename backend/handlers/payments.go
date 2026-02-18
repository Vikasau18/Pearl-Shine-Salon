package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"saloon-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentHandler struct {
	DB *pgxpool.Pool
}

func NewPaymentHandler(db *pgxpool.Pool) *PaymentHandler {
	return &PaymentHandler{DB: db}
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var req models.ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if payment exists
	var paymentID string
	var currentStatus string
	err := h.DB.QueryRow(context.Background(),
		"SELECT id, status FROM payments WHERE appointment_id = $1",
		req.AppointmentID).Scan(&paymentID, &currentStatus)

	if err != nil {
		// Create new payment
		tax := req.Amount * 0.08
		total := req.Amount + tax
		receiptNumber := fmt.Sprintf("RCP-%s-%s", time.Now().Format("20060102150405"), req.AppointmentID[:8])

		var payment models.Payment
		err = h.DB.QueryRow(context.Background(),
			`INSERT INTO payments (appointment_id, amount, discount, tax, total, method, status, receipt_number)
			 VALUES ($1, $2, 0, $3, $4, $5, 'completed', $6)
			 RETURNING id, appointment_id, amount, discount, tax, total, method, status, receipt_number, created_at, updated_at`,
			req.AppointmentID, req.Amount, tax, total, req.Method, receiptNumber,
		).Scan(&payment.ID, &payment.AppointmentID, &payment.Amount, &payment.Discount,
			&payment.Tax, &payment.Total, &payment.Method, &payment.Status,
			&payment.ReceiptNumber, &payment.CreatedAt, &payment.UpdatedAt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment"})
			return
		}

		c.JSON(http.StatusCreated, payment)
		return
	}

	// Update existing payment
	var payment models.Payment
	err = h.DB.QueryRow(context.Background(),
		`UPDATE payments SET method = $1, status = 'completed', updated_at = NOW()
		 WHERE id = $2
		 RETURNING id, appointment_id, amount, discount, tax, total, method, status, receipt_number, created_at, updated_at`,
		req.Method, paymentID,
	).Scan(&payment.ID, &payment.AppointmentID, &payment.Amount, &payment.Discount,
		&payment.Tax, &payment.Total, &payment.Method, &payment.Status,
		&payment.ReceiptNumber, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) GetSalonPayments(c *gin.Context) {
	salonID := c.Param("salon_id")

	rows, err := h.DB.Query(context.Background(),
		`SELECT p.id, p.appointment_id, p.amount, p.discount, p.tax, p.total, 
		 p.method, p.status, p.receipt_number, p.created_at, p.updated_at
		 FROM payments p
		 JOIN appointments a ON a.id = p.appointment_id
		 WHERE a.salon_id = $1
		 ORDER BY p.created_at DESC`, salonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments"})
		return
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		rows.Scan(&p.ID, &p.AppointmentID, &p.Amount, &p.Discount, &p.Tax, &p.Total,
			&p.Method, &p.Status, &p.ReceiptNumber, &p.CreatedAt, &p.UpdatedAt)
		payments = append(payments, p)
	}
	if payments == nil {
		payments = []models.Payment{}
	}

	c.JSON(http.StatusOK, payments)
}

func (h *PaymentHandler) GetReceipt(c *gin.Context) {
	paymentID := c.Param("id")

	var p models.Payment
	err := h.DB.QueryRow(context.Background(),
		`SELECT id, appointment_id, amount, discount, tax, total, method, status, receipt_number, created_at, updated_at
		 FROM payments WHERE id = $1`, paymentID,
	).Scan(&p.ID, &p.AppointmentID, &p.Amount, &p.Discount, &p.Tax, &p.Total,
		&p.Method, &p.Status, &p.ReceiptNumber, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	// Get appointment details
	var appt models.Appointment
	h.DB.QueryRow(context.Background(),
		`SELECT a.appointment_date::text, a.start_time::text, s.name, st.name, sv.name, u.name, u.email
		 FROM appointments a
		 JOIN salons s ON s.id = a.salon_id
		 JOIN staff st ON st.id = a.staff_id
		 JOIN services sv ON sv.id = a.service_id
		 JOIN users u ON u.id = a.customer_id
		 WHERE a.id = $1`, p.AppointmentID,
	).Scan(&appt.AppointmentDate, &appt.StartTime, &appt.SalonName, &appt.StaffName,
		&appt.ServiceName, &appt.CustomerName, &appt.CustomerEmail)

	c.JSON(http.StatusOK, gin.H{
		"payment":     p,
		"appointment": appt,
	})
}
