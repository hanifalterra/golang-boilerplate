package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// KafkaConsumer handles consuming messages from Kafka.
type KafkaConsumer struct {
	consumer *kafka.Consumer
	topic    string
}

// NewKafkaConsumer creates a new Kafka consumer for the specified topic.
func NewKafkaConsumer(brokers, groupID, topic string) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}
	return &KafkaConsumer{
		consumer: c,
		topic:    topic,
	}, nil
}

// Consume starts listening for Kafka messages and passes them to the handler.
func (kc *KafkaConsumer) Consume(ctx context.Context, handleFunc func(ctx context.Context, message []byte) error) {
	defer kc.consumer.Close()

	if err := kc.consumer.SubscribeTopics([]string{kc.topic}, nil); err != nil {
		log.Fatalf("failed to subscribe to topics: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer shutting down...")
			return
		default:
			ev, err := kc.consumer.ReadMessage(-1)
			if err != nil {
				log.Printf("Error reading Kafka message: %v", err)
				continue
			}
			if err := handleFunc(ctx, ev.Value); err != nil {
				log.Printf("Error handling Kafka message: %v", err)
			}
		}
	}
}
