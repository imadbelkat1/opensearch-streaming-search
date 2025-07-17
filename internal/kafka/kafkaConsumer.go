package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func NewConsumer(topic string) error {
	// to consume messages
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", GetKafkaConfig().BootstrapServers, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	log.Printf("Kafka consumer started for topic: %s", topic)

	b := make([]byte, 10e3) // 10KB max per message

	// Try to read the first message to check if there are any messages
	n, err := batch.Read(b)
	if err != nil {
		log.Printf("No messages to consume from topic: %s", topic)
		return nil
	}
	// If there is at least one message, process it and continue reading
	message := string(b[:n])
	log.Printf("Received message: %s", message)

	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		message := string(b[:n])
		log.Printf("Received message: %s", message)
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}

	log.Printf("Kafka consumer finished reading messages from topic: %s", topic)
	return nil
}
