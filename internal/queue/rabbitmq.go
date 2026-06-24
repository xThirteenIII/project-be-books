package rabbitmq

import (
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var connAttempts = 10

func AttemptConnect() (*amqp.Connection, error) {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://user:password@localhost:5672/"
	}
	var err error
	var conn *amqp.Connection
	for connAttempts > 0 {
		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ is trying to connect, attempt left: %d\n", connAttempts)
		time.Sleep(5 * time.Second)
		connAttempts--

	}
	if err != nil {
		return nil, fmt.Errorf("amqp.Dial: %v\n", err)
	}
	defer conn.Close()
	log.Printf("Connection to RabbiMQ successful.\n")
	/*
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		<-signalChan
		log.Printf("Exiting gracefully...\n")
	*/
	return conn, nil
}
