package consume

import (
	"Notifier/database"
	"Notifier/model"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ConsumeMessages() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"myTopic", "^aRegex.*[Tt]opic"}, nil)

	run := true

	for run {
		msg, err := c.ReadMessage(time.Second)
		if err == nil {
			// Mesajın içeriğini JSON olarak parse edin
			var messageData Message
			err := json.Unmarshal(msg.Value, &messageData)
			if err != nil {
				fmt.Printf("Mesaj JSON çözümlemesi hatası: %v\n", err)
			} else {
				// "Mail" ve "Code" alanlarını alın ve yazdırın
				fmt.Printf("Mail: %s, Code: %s\n", messageData.Mail, messageData.Code)

				verfication := model.Verfication{
					Email:       messageData.Mail,
					Verfication: messageData.Code,
				}

				database.Conn.Create(&verfication)

			}
		}
	}

	c.Close()
}

type Message struct {
	Mail string `json:"mail"`
	Code string `json:"code"`
}
