package services

import (
	"fmt"
	"log"
	"net/smtp"
)

type NotificationService struct {
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	SMTPFrom string
	Enabled  bool
}

func NewNotificationService(host, port, user, pass, from string) *NotificationService {
	enabled := user != "" && pass != ""
	if !enabled {
		log.Println("⚠️  Email notifications disabled (SMTP credentials not configured)")
	}
	return &NotificationService{
		SMTPHost: host,
		SMTPPort: port,
		SMTPUser: user,
		SMTPPass: pass,
		SMTPFrom: from,
		Enabled:  enabled,
	}
}

func (s *NotificationService) SendEmail(to, subject, body string) error {
	if !s.Enabled {
		log.Printf("📧 [SIMULATED] Email to %s: %s", to, subject)
		return nil
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		s.SMTPFrom, to, subject, body)

	auth := smtp.PlainAuth("", s.SMTPUser, s.SMTPPass, s.SMTPHost)
	addr := fmt.Sprintf("%s:%s", s.SMTPHost, s.SMTPPort)

	err := smtp.SendMail(addr, auth, s.SMTPFrom, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("❌ Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("✅ Email sent to %s: %s", to, subject)
	return nil
}

func (s *NotificationService) SendAppointmentConfirmation(to, customerName, salonName, date, time, serviceName string) {
	subject := "Appointment Confirmed - " + salonName
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<div style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 30px; border-radius: 12px 12px 0 0;">
				<h1 style="color: white; margin: 0;">✨ Appointment Confirmed!</h1>
			</div>
			<div style="padding: 30px; background: #f9fafb; border-radius: 0 0 12px 12px;">
				<p>Hi <strong>%s</strong>,</p>
				<p>Your appointment has been confirmed:</p>
				<div style="background: white; padding: 20px; border-radius: 8px; border-left: 4px solid #667eea;">
					<p><strong>🏪 Salon:</strong> %s</p>
					<p><strong>💇 Service:</strong> %s</p>
					<p><strong>📅 Date:</strong> %s</p>
					<p><strong>🕐 Time:</strong> %s</p>
				</div>
				<p style="margin-top: 20px; color: #6b7280;">See you there!</p>
			</div>
		</body>
		</html>
	`, customerName, salonName, serviceName, date, time)

	go s.SendEmail(to, subject, body)
}

func (s *NotificationService) SendAppointmentReminder(to, customerName, salonName, date, time string) {
	subject := "Reminder: Your appointment at " + salonName
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<div style="background: linear-gradient(135deg, #f093fb 0%%, #f5576c 100%%); padding: 30px; border-radius: 12px 12px 0 0;">
				<h1 style="color: white; margin: 0;">⏰ Appointment Reminder</h1>
			</div>
			<div style="padding: 30px; background: #f9fafb; border-radius: 0 0 12px 12px;">
				<p>Hi <strong>%s</strong>,</p>
				<p>This is a friendly reminder of your upcoming appointment:</p>
				<div style="background: white; padding: 20px; border-radius: 8px; border-left: 4px solid #f5576c;">
					<p><strong>🏪 Salon:</strong> %s</p>
					<p><strong>📅 Date:</strong> %s</p>
					<p><strong>🕐 Time:</strong> %s</p>
				</div>
			</div>
		</body>
		</html>
	`, customerName, salonName, date, time)

	go s.SendEmail(to, subject, body)
}
