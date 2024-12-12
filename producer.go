package main

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

// UserEvent represents a user activity
type UserEvent struct {
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}

// OrderEvent represents an order in the system
type OrderEvent struct {
	OrderID    string  `json:"order_id"`
	CustomerID string  `json:"customer_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
}

func main() {
	topic := "test-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial producer:", err)
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatal("failed to set write deadline:", err)
	}

	// Example 1: Simple string message
	simpleMessage := kafka.Message{
		Value: []byte("Example 1: Hello, this is a test message from gononet!"),
	}

	// Example 2: JSON message with user event
	userEvent := UserEvent{
		UserID:    "Example 2",
		Action:    "micro-managing",
		Timestamp: time.Now(),
	}
	userEventJSON, _ := json.Marshal(userEvent)
	userMessage := kafka.Message{
		Key:   []byte("user_event"),
		Value: userEventJSON,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte("user_activity")},
			{Key: "version", Value: []byte("2.9")},
		},
		Time: time.Time{},
	}

	// Example 3: Order event with different headers
	orderEvent := OrderEvent{
		OrderID:    "Example 3",
		CustomerID: "your_customer_id",
		Amount:     69.69,
		Status:     "Active",
	}
	orderEventJSON, _ := json.Marshal(orderEvent)
	orderMessage := kafka.Message{
		Key:   []byte("order_event"),
		Value: orderEventJSON,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte("order_processing")},
			{Key: "priority", Value: []byte("high")},
		},
	}

	// Write all messages in batch
	_, err = conn.WriteMessages(
		simpleMessage,
		userMessage,
		orderMessage,
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	log.Println("Successfully wrote messages to Kafka")

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
