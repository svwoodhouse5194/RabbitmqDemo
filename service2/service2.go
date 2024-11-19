package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	Name    string
	Body    string
	Time    int64
	History []string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func sendMessage(message Message) {
	// Connects to the RabbitMQ server
	host, ok := os.LookupEnv("RABBITMQ_SVC_NAME")
	if !ok {
		host = "localhost"
	}

	conn, err := amqp.Dial("amqp://guest:guest@" + host + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declares the queue
	q, err := ch.QueueDeclare(
		"pythonQueue", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Adds new message to the message.history array
	message.History = append(message.History, "Message was received to Service 2")
	message.History = append(message.History, "Message was sent to Service 3")

	// Changes the body of the message to reflect current state
	message.Body = "Sent from Go Service 2"

	// Converts message to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(jsonData),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent message to Service 3")
}

func main() {
	// Connects to RabbitMQ server
	host, ok := os.LookupEnv("RABBITMQ_SVC_NAME")
	if !ok {
		host = "localhost"
	}

	conn, err := amqp.Dial("amqp://guest:guest@" + host + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declares the queue
	q, err := ch.QueueDeclare(
		"goQueue", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Consumes the messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			// Decodes the message
			var message Message
			dec := gob.NewDecoder(bytes.NewReader(d.Body))

			err = dec.Decode(&message)
			if err != nil {
				fmt.Println("Error decoding:", err)
				return
			}
			log.Printf("Received a message from Service 1")

			// Sends message to next service
			sendMessage(message)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
