package email

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"os"
	"log"
)

func RunEmailWorker(rabbitMQURL string) error {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare("email_queue", true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	sendgridAPIKey := os.Getenv("SENDGRID_API_KEY")
	if sendgridAPIKey == "" {
		log.Println("SENDGRID_API_KEY is not set")
		return nil
	}

	for msg := range msgs {
		var email struct {
			To       string `json:"to"`
			UserName string `json:"user_name"`
			Body     string `json:"body"`
		}
		if err := json.Unmarshal(msg.Body, &email); err != nil {
			log.Printf("Failed to unmarshal email message: %v", err)
			continue
		}

		from := mail.NewEmail("iSymptom", "noreply@isymptom.vercel.app")
		to := mail.NewEmail(email.UserName, email.To)
		subject := "Request to reset password"
		htmlContent := `<div style=\"font-family: Arial, sans-serif; padding: 20px;\">` +
			`<h2>Xin chào ` + email.UserName + `,</h2>` +
			`<p>Bạn vừa yêu cầu đặt lại mật khẩu cho tài khoản iSymptom.</p>` +
			`<a href=\"` + email.Body + `\" style=\"display: inline-block; padding: 10px 20px; background-color: #007bff; color: #fff; text-decoration: none; border-radius: 5px;\">Đặt lại mật khẩu</a>` +
			`<p>Nếu bạn không yêu cầu, hãy bỏ qua email này.</p>` +
			`</div>`

		message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
		client := sendgrid.NewSendClient(sendgridAPIKey)
		_, err := client.Send(message)
		if err != nil {
			log.Printf("Failed to send email via SendGrid: %v", err)
		}
	}
	return nil
}