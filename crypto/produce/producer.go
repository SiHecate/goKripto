package produce

import (
	model "cryptoApp/Model"
	"cryptoApp/database"
	"cryptoApp/helpers"
	"log"

	"encoding/json"

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

	message := struct {
		Code string `json:"code"`
		Mail string `json:"mail"`
	}{
		Code: valid_code,
		Mail: email,
	}

	// JSON nesnesini Kafka'ya göndermek için marshal edin.
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Marshal error: %v\n", err)
		return err
	}

	// JSON nesnesini Kafka'ya gönderin.
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(email),
		Value:          jsonMessage,
	}, nil)
	if err != nil {
		log.Printf("Producer error: %v\n", err)
		return err
	}

	p.Flush(15 * 1000)

	return nil
}
