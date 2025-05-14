package models

import "github.com/streadway/amqp"

type RabbitMQClient interface {
	CreateQueue(queueName string) error
	Publish(queueName string, message []byte) error
	Consume(queueName string) (<-chan amqp.Delivery, error)
	PublishEmailMessage(email, userName, resetLink string) error
	ConsumeEmailMessages(handler func(email, userName, resetLink string) error) error
	Close() error
	IsConnected() bool
}
