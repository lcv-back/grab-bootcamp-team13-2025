package mocks

import (
	"context"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

type RabbitMQClient struct {
	mock.Mock
}

func (m *RabbitMQClient) PublishEmailMessage(email, userName, resetLink string) error {
	args := m.Called(email, userName, resetLink)
	return args.Error(0)
}

func (m *RabbitMQClient) ConsumeEmailMessages(handler func(email, userName, resetLink string) error) error {
	args := m.Called(handler)
	return args.Error(0)
}

func (m *RabbitMQClient) Close() error {
	m.Called()
	return nil
}

func (m *RabbitMQClient) CreateQueue(queueName string) error {
	args := m.Called(queueName)
	return args.Error(0)
}

func (m *RabbitMQClient) Publish(queueName string, message []byte) error {
	args := m.Called(queueName, message)
	return args.Error(0)
}

func (m *RabbitMQClient) Consume(queueName string) (<-chan amqp.Delivery, error) {
	args := m.Called(queueName)
	return args.Get(0).(<-chan amqp.Delivery), args.Error(1)
}

func (m *RabbitMQClient) IsConnected() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *RabbitMQClient) SendResetPasswordEmail(ctx context.Context, fullname, email, resetLink string) error {
	args := m.Called(ctx, fullname, email, resetLink)
	return args.Error(0)
}
