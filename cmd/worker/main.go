package main

import (
	"context"
	"log"

	"golang-boilerplate/internal/pkg/connections/kafka"
	"golang-boilerplate/internal/pkg/services"
	"golang-boilerplate/internal/worker/controllers"
)

func main() {
	ctx := context.Background()

	// Kafka configuration.
	brokers := "localhost:9092"
	groupID := "transaction-consumer-group"
	topic := "transaction"

	// Initialize Kafka consumer.
	consumer, err := kafka.NewKafkaConsumer(brokers, groupID, topic)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka consumer: %v", err)
	}

	// Initialize service and controller.
	service := services.NewTransactionService()
	controller := controllers.NewTransactionController(service)

	// Start consuming messages.
	log.Println("Starting Kafka consumer...")
	consumer.Consume(ctx, controller.HandleMessage)
}
