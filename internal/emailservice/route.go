package emailservice

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ishola-faazele/taskflow/pkg/utils/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

func RegisterRoutes(conn *amqp.Connection) {
	logger := logger.NewStdLogger()
	es := NewEmailService(DefaultEmailConfig())
	c := NewEmailConsumer(es)
	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("FAILED_TO_OPEN_CHANNEL")
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"email_queue", // queue
		"",            // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		log.Panicf("FAILED_TO_REGISTER_CONSUMER")
	}
	logger.Info("Listening to email_queue")

	forever := make(chan struct{})
	var msg EmailMessage
	go func() {
		for d := range msgs {
			logger.Info("New Email Message Received.")
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				return
			}
			err = c.HandleEmailMessage(msg)
			if err != nil {
				return
			}
			err = d.Ack(false)
			if err != nil {
				errMsg := fmt.Sprintf("FAILED_TO_ACKNOWLEDGE_MESSAGE: %v", err)
				logger.Error(errMsg)
				break
			}
			logger.Info("Email Message Acknowledged.")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
