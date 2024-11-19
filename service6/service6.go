package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	/// Connects to the RabbitMQ server
	host, ok := os.LookupEnv("RABBITMQ_SVC_NAME")
	if !ok {
		host = "localhost"
	}
	conn, err := amqp.Dial("amqp://guest:guest@" + host + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Opens the channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declares the exchange
	err = ch.ExchangeDeclare(
		"siblingExchange", // name
		"fanout",          // type
		false,             // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// Declares the queue
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Binds the Queue
	err = ch.QueueBind(
		q.Name,            // queue name
		"",                // routing key
		"siblingExchange", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	// Gets the messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		true,   // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message from Service 4")
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
