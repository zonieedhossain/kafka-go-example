# Kafka Example: Consumer and Producer

This repository contains examples of a Kafka consumer and producer implemented in Go using the [segmentio/kafka-go](https://github.com/segmentio/kafka-go) library. The examples demonstrate:

- Setting up a Kafka consumer to read messages from a topic.
- Setting up a Kafka producer to write messages to a topic.
- Handling different types of events and formats, including simple strings and JSON objects.
- Best practices for managing Kafka connections, message deletion, and error handling.

---

## Prerequisites

1. **Kafka Cluster**: A running Kafka instance accessible on `localhost:9092`.
2. **Go Environment**: Installed Go compiler (1.16 or higher is recommended).
3. **Dependencies**: The `kafka-go` package should be installed. You can add it to your project with:
```bash
go get github.com/segmentio/kafka-go
```

---

## Examples Overview

### 1. Kafka Consumer Example

The consumer reads messages from a Kafka topic (`test-topic`) and processes them sequentially. It also demonstrates:
- Setting a read deadline.
- Handling end-of-topic and temporary errors gracefully.
- Writing tombstone messages to simulate message deletion.

#### Key Functions
- `deleteMessage(writer *kafka.Writer, originalMsg kafka.Message)`: Sends a tombstone message to delete a Kafka message.

### 2. Kafka Producer Example

The producer writes messages to the same topic (`test-topic`). It showcases:
- Sending simple text messages.
- Sending JSON-encoded events for structured data (e.g., `UserEvent` and `OrderEvent`).
- Adding headers for metadata.

---

## How to Run

### Step 1: Start Kafka
Ensure your Kafka server is running locally at `localhost:9092`. You can use Docker to quickly set up a Kafka instance:

```bash
docker-compose up -d
```

Or start Kafka manually.

### Step 2: Run the Consumer
Navigate to the directory containing the consumer code and run:

```bash
go run consumer.go
```

The consumer will read messages from the `test-topic` and process them. If the topic is empty, it will gracefully exit.

### Step 3: Run the Producer
Navigate to the directory containing the producer code and run:

```bash
go run producer.go
```

The producer will send the following messages to the `test-topic`:
- A simple string message.
- A JSON-encoded `UserEvent`.
- A JSON-encoded `OrderEvent`.

### Expected Output

#### Consumer
For each message read, the consumer prints:
```text
Message <count>: <message content>
```
If it reaches the end of the topic, it prints:
```text
Reached end of topic after <count> messages
```

#### Producer
After sending messages, the producer logs:
```text
Successfully wrote messages to Kafka
```
---
## Code Structure

### Consumer Code Highlights
- **Connection**: Uses `kafka.DialLeader` to connect to the Kafka leader for the specified topic and partition.
- **Message Processing**: Reads messages in a loop and handles errors appropriately.
- **Message Deletion**: Sends tombstone messages to simulate deletion.

### Producer Code Highlights
- **Connection**: Uses `kafka.DialLeader` to establish a producer connection.
- **Event Handling**: Demonstrates structured events like `UserEvent` and `OrderEvent` with custom headers.
- **Batch Writing**: Writes multiple messages in a single batch for efficiency.
