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
		log.Fatalf("error creating environment: %s", err)
	}

	streamName := "mixing"
	producer, err := env.NewProducer(streamName, nil)
	if err != nil {
		log.Fatalf("error creating producer: %s", err)
	}

	//optional publish confirmation channel
	chPublishConfirm := producer.NotifyPublishConfirmation()
	handlePublishConfirm(chPublishConfirm)

	expiryTime := time.Now().Local().Add(time.Second * time.Duration(600))

	// send basic messages
	var messages []message.StreamMessage
	for i := 1; i < 100; i++ {
		data := fmt.Sprintf("message %v: with body only", i)
		msg := amqp.NewMessage([]byte(data))
		messages = append(messages, msg)
	}

	// send messages with body and properties
	for i := 1; i < 100; i++ {
		data := fmt.Sprintf("message %v: body and properties", i)
		msg := amqp.NewMessage([]byte(data))
		msg.Properties = &amqp.MessageProperties{
			MessageID:          i,
			UserID:             []byte{'1'},
			To:                 "0",
			Subject:            "second message",
			ReplyTo:            "0",
			CorrelationID:      i * 2,
			ContentType:        "text/string",
			ContentEncoding:    "utf-8",
			AbsoluteExpiryTime: expiryTime,
			CreationTime:       time.Now(),
			GroupID:            "groupID",
			GroupSequence:      8,
			ReplyToGroupID:     "ReplyToGroupID",
		}
		messages = append(messages, msg)
	}

	// send messages with body, properties, and app properties
	for i := 1; i < 100; i++ {
		data := fmt.Sprintf("message %v: with body, properties and app properties", i)
		msg := amqp.NewMessage([]byte(data))
		msg.Properties = &amqp.MessageProperties{
			MessageID:          i * 4,
			UserID:             []byte{'1'},
			To:                 "0",
			Subject:            "third message",
			ReplyTo:            "0",
			CorrelationID:      i * 5,
			ContentType:        "text/string",
			ContentEncoding:    "utf-8",
			AbsoluteExpiryTime: expiryTime,
			CreationTime:       time.Now(),
			GroupID:            "groupID",
			GroupSequence:      9,
			ReplyToGroupID:     "ReplyToGroupID",
		}
		msg.ApplicationProperties = map[string]interface{}{
			"key_string":   "values",
			"key2_int":     "1111",
			"key2_decimal": 10.000,
		}
		messages = append(messages, msg)
	}

	err = producer.BatchSend(messages)

	if err != nil {
		log.Fatalf("error sending message: %s", err)
	}

	// sleep to show confirmed messages
	time.Sleep(1 * time.Second)
	err = producer.Close()
	if err != nil {
		log.Fatalf("error closing producer: %s", err)
	}
}
