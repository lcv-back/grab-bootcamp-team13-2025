package main

import (
	"grab-bootcamp-be-team13-2025/internal/config"
	"grab-bootcamp-be-team13-2025/pkg/utils/email"
	"grab-bootcamp-be-team13-2025/pkg/utils/rabbitmq"
	"log"
	"os"
)

func main() {
	// load configuration
	ctg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatal("failed to load config", err)
	}

	// initialize email sender
	emailSender := email.NewEmailSender(ctg)

	// connect to rabbitmq
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	rabbitClient, err := rabbitmq.NewRabbitMQClient(rabbitmqURL)
	if err != nil {
		log.Fatal("failed to connect to rabbitmq", err)
	}

	defer rabbitClient.Close()

	// handle messages from queue
	err = rabbitClient.ConsumeEmailMessages(func(email, userName, resetLink string) error {
		return emailSender.SendResetPasswordEmail(email, resetLink)
	})
	if err != nil {
		log.Fatal("failed to consume messages from queue", err)
	}

	// keep worker running
	select {}
}
