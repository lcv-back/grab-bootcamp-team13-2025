// CẢNH BÁO: File này không còn được khuyến khích sử dụng cho gửi email reset password.
// Hãy sử dụng internal/infrastructure/email/email.service.go và worker.go để gửi email HTML chuyên nghiệp.

package email

import (
	"fmt"
	"grab-bootcamp-be-team13-2025/internal/config"
	"grab-bootcamp-be-team13-2025/pkg/utils/rabbitmq"
	"net/smtp"
)

type EmailSender struct {
	config *config.Config
	auth   smtp.Auth
}

func NewEmailSender(config *config.Config) *EmailSender {
	auth := smtp.PlainAuth("", config.Email.Username, config.Email.Password, config.Email.Host)
	return &EmailSender{
		config: config,
		auth:   auth,
	}
}

func (e *EmailSender) SendResetPasswordEmail(toEmail, resetLink string) error {
	//from := e.config.Email.Username
	from := "no-reply@isymptom.vercel.app"
	to := []string{toEmail}
	subject := "Request to reset password"
	body := fmt.Sprintf("Click <a href=\"%s\">here</a> to reset your password.", resetLink)

	msg := []byte(
		"To: " + toEmail + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	smtpAddr := fmt.Sprintf("%s:%d", e.config.Email.Host, e.config.Email.Port)
	err := smtp.SendMail(smtpAddr, e.auth, from, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func RunEmailWorker(host string, port int, username, password, rabbitmqURL string) error {
	// Create config
	cfg := &config.Config{}
	cfg.Email.Host = host
	cfg.Email.Port = port
	cfg.Email.Username = username
	cfg.Email.Password = password

	// Initialize email sender
	emailSender := NewEmailSender(cfg)

	// Connect to rabbitmq
	rabbitClient, err := rabbitmq.NewRabbitMQClient(rabbitmqURL)
	if err != nil {
		return fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	defer rabbitClient.Close()

	// Handle messages from queue
	err = rabbitClient.ConsumeEmailMessages(func(email, userName, resetLink string) error {
		return emailSender.SendResetPasswordEmail(email, resetLink)
	})
	if err != nil {
		return fmt.Errorf("failed to consume messages from queue: %w", err)
	}

	// Keep worker running
	select {}
}
