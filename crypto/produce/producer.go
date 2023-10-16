package produce

import (
	model "cryptoApp/Model"
	"cryptoApp/database"
	"cryptoApp/helpers"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
)

func VerficationCode(c *fiber.Ctx) error {
	issuer, err := helpers.GetIssuer(c)
	if err != nil {
		return err
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka",
		"linger.ms":         10,
		"batch.size":        1000,
	})
	if err != nil {
		log.Printf("Producer creation error: %v\n", err)
		return err
	}
	defer p.Close()

	topic := "myTopic"

	email, err := model.GetMailByIssuer(database.Conn, issuer)
	if err != nil {
		return err
	}

	valid_code, err := model.GetVerficationByMail(database.Conn, email)
	if err != nil {
		return err
	}

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(valid_code),
	}, nil)
	if err != nil {
		log.Printf("Producer error: %v\n", err)
		return err
	}

	p.Flush(15 * 1000)

	return nil
}
