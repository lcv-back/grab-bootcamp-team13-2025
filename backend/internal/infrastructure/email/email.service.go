package email

import (
    "context"
    "fmt"
    "github.com/sendgrid/sendgrid-go"
    "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService interface {
    SendResetPasswordEmail(ctx context.Context, userName, toEmail, resetLink string) error
}

type emailService struct {
    apiKey string
}

func NewEmailService(apiKey string) EmailService {
    return &emailService{apiKey: apiKey}
}

func (s *emailService) SendResetPasswordEmail(ctx context.Context, userName, toEmail, resetLink string) error {
    from := mail.NewEmail("iSymptom", "noreply@isymptom.id.vn")
    to := mail.NewEmail(userName, toEmail)
    subject := "Request to reset password"
    htmlContent := fmt.Sprintf(`
    <div style="font-family: Arial, sans-serif; padding: 20px;">
        <h2>Hello %s,</h2>
        <p>You have requested to reset your password for your iSymptom account.</p>
        <a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #007bff; color: #fff; text-decoration: none; border-radius: 5px;">Reset password</a>
        <p>If you did not request this, please ignore this email.</p>
    </div>
    `, userName, resetLink)
    message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
    client := sendgrid.NewSendClient(s.apiKey)
    _, err := client.Send(message)
    return err
}