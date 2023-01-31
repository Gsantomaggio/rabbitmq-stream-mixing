package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

func main() {
	fmt.Println("consumer")

	streamName := "mixing"

	env, err := stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			SetHost("localhost").
			SetPort(5552).
			SetUser("guest").
			SetPassword("guest"),
	)

	if err != nil {
		log.Fatalf("error creating environment: %s", err)
	}

	handleMessages := func(consumerContext stream.ConsumerContext, message *amqp.Message) {
		fmt.Printf("consumer name: %s, text: %s \n ", consumerContext.Consumer.GetName(), message.Data)
	}

	consumer, err := env.NewConsumer(
		streamName,
		handleMessages,
		stream.NewConsumerOptions().
			SetConsumerName("my_consumer").                  // set a consumer name
			SetOffset(stream.OffsetSpecification{}.First()). // start consuming from the beginning
			SetCRCCheck(false))                              // Disable crc control, increase the performances

	if err != nil {
		log.Fatalf("error creating environment: %s", err)
	}

	time.Sleep(2000 * time.Millisecond)
	err = consumer.Close()
	if err != nil {
		log.Fatalf("error closing consumer: %s", err)
	}
	err = env.Close()
	if err != nil {
		log.Fatalf("error closing environment: %s", err)
	}
}
