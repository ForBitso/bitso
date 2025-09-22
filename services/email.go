package services

import (
	"fmt"
	"log"

	"go-shop/config"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

func (es *EmailService) SendOTPEmail(email, otp string) error {
	subject := "Your OTP Code"
	body := fmt.Sprintf(`
		<h2>Your OTP Code</h2>
		<p>Your OTP code is: <strong>%s</strong></p>
		<p>This code will expire in %d minutes.</p>
		<p>If you didn't request this code, please ignore this email.</p>
	`, otp, es.config.OTP.ExpireMinutes)

	return es.sendEmail(email, subject, body)
}

func (es *EmailService) SendPasswordResetEmail(email, otp string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
		<h2>Password Reset Request</h2>
		<p>You requested to reset your password.</p>
		<p>Your OTP code is: <strong>%s</strong></p>
		<p>This code will expire in %d minutes.</p>
		<p>If you didn't request this, please ignore this email.</p>
	`, otp, es.config.OTP.ExpireMinutes)

	return es.sendEmail(email, subject, body)
}

func (es *EmailService) SendWelcomeEmail(email, firstName string) error {
	subject := "Welcome to Go Shop!"
	body := fmt.Sprintf(`
		<h2>Welcome %s!</h2>
		<p>Thank you for registering with Go Shop.</p>
		<p>Your account has been successfully created and activated.</p>
		<p>Happy shopping!</p>
	`, firstName)

	return es.sendEmail(email, subject, body)
}

func (es *EmailService) sendEmail(to, subject, body string) error {
	if es.config.Email.SMTPUsername == "" || es.config.Email.SMTPPassword == "" {
		log.Printf("Email not configured, would send to %s: %s", to, subject)
		return fmt.Errorf("email service not configured - missing SMTP credentials")
	}

	log.Printf("Attempting to send email to %s via %s:%d", to, es.config.Email.SMTPHost, es.config.Email.SMTPPort)

	m := gomail.NewMessage()
	m.SetHeader("From", es.config.Email.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		es.config.Email.SMTPHost,
		es.config.Email.SMTPPort,
		es.config.Email.SMTPUsername,
		es.config.Email.SMTPPassword,
	)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}
