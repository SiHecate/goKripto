package consume

import (
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

	type message struct {
		Mail string `json:"mail"`
		Code string `json:"code"`
	}

	for run {
		msg, err := c.ReadMessage(time.Second)
		if err == nil {
			// Mesajın içeriğini JSON olarak parse edin
			var messageData message
			err := json.Unmarshal(msg.Value, &messageData)
			if err != nil {
				fmt.Printf("Mesaj JSON çözümlemesi hatası: %v\n", err)
			} else {
				// "Mail" ve "Code" alanlarını alın ve yazdırın
				fmt.Printf("Mail: %s, Code: %s\n", messageData.Mail, messageData.Code)
			}
		}
	}

	c.Close()
}
