package amqp

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Connects with AMQP Server and initialize queues
func InitAMQP() *amqp.Connection {
	// Connect to RabbitMQ
	user := os.Getenv("RABBITMQ_USER")
	pass := os.Getenv("RABBITMQ_PASS")
	addr := fmt.Sprintf("amqp://%s:%s@localhost:5673/", user, pass)

	conn, err := amqp.Dial(addr)
	failOnError(err, "FAILED_TO_CONNECT_TO_RABBITMQ")

	ch, err := conn.Channel()
	failOnError(err, "FAILED_TO_OPEN_CHANNEL")

	_, err = ch.QueueDeclare(
		"email_queue", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "FAILED_TO_DECLARE_QUEUE")
	return conn
}
