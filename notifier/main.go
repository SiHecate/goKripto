package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} - ${method} ${path}\n${body}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
	}))

	app.Post("/send-message", func(c *fiber.Ctx) error {
		data := c.FormValue("message")
		p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka"})
		if err != nil {
			log.Printf("Producer creation error: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Producer creation error",
			})
		}
		defer p.Close()

		topic := "myTopic"
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(data),
		}, nil)

		p.Flush(15 * 1000)

		return c.JSON(fiber.Map{
			"message": "Veri Kafka'ya gönderildi",
		})
	})

	log.Fatal(app.Listen(":8082"))
}
