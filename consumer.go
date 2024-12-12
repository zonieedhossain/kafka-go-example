package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func main() {
	topic := "test-topic"
	partition := 0

	// Create consumer connection
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial consumer:", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing consumer connection: %v", err)
		}
	}()

	// Important: Set read deadline
	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		log.Printf("Error setting read deadline: %v", err)
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP("localhost:9092"),
		Topic:        topic,
		Balancer:     nil,
		RequiredAcks: kafka.RequireAll,
	}
	defer func() {
		if err := writer.Close(); err != nil {
			log.Printf("Error closing writer: %v", err)
		}
	}()

	messageCount := 0
	for {
		message, err := conn.ReadMessage(10e3)

		// Proper error handling
		if err != nil {
			if err.Error() == "EOF" {
				log.Printf("Reached end of topic after %d messages", messageCount)
				break
			}

			// Check for timeout
			if err.(kafka.Error).Temporary() {
				log.Printf("No more messages available. Total processed: %d", messageCount)
				break
			}

			log.Printf("Error reading message: %v", err)
			break
		}

		messageCount++
		fmt.Printf("Message %d: %s\n", messageCount, string(message.Value))

		if err = deleteMessage(writer, message); err != nil {
			log.Printf("Failed to delete message: %v", err)
			break
		}
	}

	log.Printf("Processing complete. Total messages: %d", messageCount)
}
func deleteMessage(writer *kafka.Writer, originalMsg kafka.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tombstone := kafka.Message{
		Key:       originalMsg.Key,
		Value:     nil,
		Partition: originalMsg.Partition,
	}

	err := writer.WriteMessages(ctx, tombstone)
	if err != nil {
		return fmt.Errorf("error writing tombstone message: %w", err)
	}

	return nil
}
