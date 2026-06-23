package rabbitmq

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func AttemptConnect() error {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://user:password@localhost:5672/"
	}
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %v\n", err)
	}
	defer conn.Close()
	log.Printf("Connection to RabbiMQ successful.\n")
	/*
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		<-signalChan
		log.Printf("Exiting gracefully...\n")
	*/
	return nil
}
