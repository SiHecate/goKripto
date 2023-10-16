package main

import (
	"Notifier/database"
	"log"
	"math/rand"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func GenerateCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var codes [6]byte
	for i := 0; i < 6; i++ {
		codes[i] = uint8(48 + r.Intn(10))
	}

	return string(codes[:])
}

func main() {
	database.Connect()
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka"})
	if err != nil {
		log.Printf("Producer creation error: %v\n", err)
		return
	}
	defer p.Close()

	topic := "myTopic"

	// Generate a random code
	code := GenerateCode()

	// Send the code to Kafka
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(code),
	}, nil)

	p.Flush(15 * 1000)
}
