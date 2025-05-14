package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = channel.QueueDeclare(
		"email_queue",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQClient{
		conn:    conn,
		channel: channel,
	}, nil
}

// CreateQueue tạo một queue trong RabbitMQ
func (c *RabbitMQClient) CreateQueue(queueName string) error {
	_, err := c.channel.QueueDeclare(
		queueName, // Tên queue
		true,      // Durable
		false,     // Auto-delete
		false,     // Exclusive
		false,     // No-wait
		nil,       // Args
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %v", queueName, err)
	}
	return nil
}

// Publish gửi một message tới queue
func (c *RabbitMQClient) Publish(queueName string, message []byte) error {
	err := c.channel.Publish(
		"",        // Exchange
		queueName, // Routing key (tên queue)
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to queue %s: %v", queueName, err)
	}
	return nil
}

// Consume trả về một channel để nhận message từ queue
func (c *RabbitMQClient) Consume(queueName string) (<-chan amqp.Delivery, error) {
	messages, err := c.channel.Consume(
		queueName, // Tên queue
		"",        // Consumer tag
		true,      // Auto-ack
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages from queue %s: %v", queueName, err)
	}
	return messages, nil
}

// Fix the method name from PublicEmailMessage to PublishEmailMessage
func (r *RabbitMQClient) PublishEmailMessage(email, userName, resetLink string) error {
	msg := struct {
		Email     string `json:"email"`
		UserName  string `json:"user_name"`
		ResetLink string `json:"reset_link"`
	}{
		Email:     email,
		UserName:  userName,
		ResetLink: resetLink,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = r.channel.Publish(
		"",
		"email_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQClient) ConsumeEmailMessages(handler func(email, userName, resetLink string) error) error {
	msgs, err := r.channel.Consume(
		"email_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var emailMsg struct {
				Email     string `json:"email"`
				UserName  string `json:"user_name"`
				ResetLink string `json:"reset_link"`
			}
			if err := json.Unmarshal(msg.Body, &emailMsg); err != nil {
				continue
			}
			handler(emailMsg.Email, emailMsg.UserName, emailMsg.ResetLink)
		}
	}()

	return nil
}

func (r *RabbitMQClient) Close() error {
	var firstErr error

	if r.channel != nil {
		if err := r.channel.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// Add this method to the RabbitMQClient struct
func (r *RabbitMQClient) IsConnected() bool {
	return r.conn != nil && !r.conn.IsClosed()
}
