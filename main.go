package main

import (
	"fmt"
	"gokripto/database"
	router "gokripto/router"
	websocket "gokripto/router"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "gokripto/docs"

	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// @title           Go Crypto
// @version         1.0
// @description     Crypto currency app.
// @contact.name   API Support
// @contact.url    https://github.com/SiHecate
// @host      localhost:8080
// @securityDefinitions.basic  BasicAuth
func main() {
	database.Connect()
	app := fiber.New()

	go func() {
		consumer()
	}()

	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} - ${method} ${path}\n${body}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
	}))

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	router.Setup(app)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		websocket.StartWebSocket(app)
		wg.Done()
	}()
	wg.Wait()

	app.Listen(":8080")
}

func consumer() {
	topic := "comments"
	worker, err := connectConsumer([]string{"kafka:9092"})
	if err != nil {
		panic(err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	fmt.Println("Consumer started ")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	msgCount := 0
	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Received message Count %d: | Topic(%s) | Message(%s) \n", msgCount, string(msg.Topic), string(msg.Value))
			case <-sigchan:
				fmt.Println("Interrupt is detected")
				close(doneCh)
				return
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCount, "messages")

	if err := worker.Close(); err != nil {
		panic(err)
	}
}

func connectConsumer(brokersUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
