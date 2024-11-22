package main

import (
	"bytes"
	"context"
	"encoding/gob"
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

func main() {
	// Connects to the RabbitMQ server
	host, ok := os.LookupEnv("RABBITMQ_SVC_NAME")
	if !ok {
		host = "localhost"
	}

	user, ok := os.LookupEnv("USER")
	if !ok {
		host = "guest"
	}

	pass, ok := os.LookupEnv("PASS")
	if !ok {
		host = "guest"
	}

	conn, err := amqp.Dial("amqp://" + user + ":" + pass + "@" + host + ":5672/")

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Creates the message
	m := Message{Name: "Jane Doe", Body: "Sent from Go Service 1", Time: time.Now().Unix(), History: []string{"Message was created in Service 1"}}

	// Create a buffer to hold the serialized data
	var buf bytes.Buffer

	// Create a new encoder and encode the struct
	enc := gob.NewEncoder(&buf)
	er := enc.Encode(m)
	if er != nil {
		fmt.Println("Error:", er)
		return
	}

	// Sends the message to the queue
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        buf.Bytes(),
		})

	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent message to Service 2")
}
