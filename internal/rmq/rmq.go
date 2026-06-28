package rmq

import (
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const connAttempts = 10

func AttemptConnect() (*amqp.Connection, error) {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://user:password@localhost:5672/"
	}
	var err error
	var conn *amqp.Connection
	for attemptsLeft := connAttempts; attemptsLeft > 0; attemptsLeft-- {
		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ is trying to connect, attempts left: %d\n", attemptsLeft)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("amqp.Dial: %v\n", err)
	}
	log.Printf("Connection to RabbitMQ successful.\n")
	return conn, nil
}
