package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/message"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

func handlePublishConfirm(confirms stream.ChannelPublishConfirm) {
	go func() {
		for confirmed := range confirms {
			for _, msg := range confirmed {
				if msg.IsConfirmed() {
					fmt.Printf("message %s stored \n  ", msg.GetMessage().GetData())
				} else {
					fmt.Printf("message %s failed \n  ", msg.GetMessage().GetData())
				}

			}
		}
	}()
}

func main() {
	fmt.Println("producer")
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

	streamName := "mixing"
	err = env.DeclareStream(streamName,
		stream.NewStreamOptions().
			SetMaxLengthBytes(stream.ByteCapacity{}.GB(2)))

	if err != nil {
		log.Fatalf("error declaring stream")
	}

	producer, err := env.NewProducer(streamName, nil)
	if err != nil {
		log.Fatalf("error creating producer")
	}

	//optional publish confirmation channel
	chPublishConfirm := producer.NotifyPublishConfirmation()
	handlePublishConfirm(chPublishConfirm)

	var message message.StreamMessage
	message = amqp.NewMessage([]byte("a message send by go-stream-client"))
	err = producer.Send(message)
	if err != nil {
		log.Fatalf("error sending message")
	}

	// sleep to show confirmed messages
	time.Sleep(1 * time.Second)
	err = producer.Close()
	if err != nil {
		log.Fatalf("error closing producer")
	}
}
