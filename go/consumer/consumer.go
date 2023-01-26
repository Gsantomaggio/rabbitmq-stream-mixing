package main

import (
	"fmt"
	"log"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

func consumerClose(channelClose stream.ChannelClose) {
	event := <-channelClose
	fmt.Printf("Consumer: %s closed on the stream: %s, reason: %s \n", event.Name, event.StreamName, event.Reason)
}

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
		log.Fatalf("error creating environment")
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
		log.Fatalf("error creating environment")
	}

	channelClose := consumer.NotifyClose()
	// channelClose receives all the closing events, here you can handle the
	// client reconnection or just log
	defer consumerClose(channelClose)

	err = env.DeleteStream(streamName)
	if err != nil {
		log.Fatalf("error deleting stream")
	}
	err = env.Close()
	if err != nil {
		log.Fatalf("error closing environment")
	}
}
