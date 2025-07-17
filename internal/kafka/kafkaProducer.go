package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func NewItemProducer(topic string, ids []int) error {
	// to produce messages
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", GetKafkaConfig().BootstrapServers, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// Create messages for all IDs
	messages := make([]kafka.Message, len(ids))
	for i, id := range ids {
		messages[i] = kafka.Message{Value: fmt.Appendf(nil, "%d", id)}
	}

	_, err = conn.WriteMessages(messages...)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	log.Printf("Kafka producer sent %d messages with IDs: %v", len(ids), ids)
	return nil
}

func NewUserIDProducer(topic string, ids []string) error {
	// to produce messages
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", GetKafkaConfig().BootstrapServers, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// Create messages for all user IDs
	messages := make([]kafka.Message, len(ids))
	for i, id := range ids {
		messages[i] = kafka.Message{Value: []byte(id)}
	}

	_, err = conn.WriteMessages(messages...)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	log.Printf("Kafka producer sent %d user ID messages: %v", len(ids), ids)
	return nil
}
