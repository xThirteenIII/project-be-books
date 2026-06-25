package rmq

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartReviewConsumer(ch *amqp.Channel, queueName string) error {
	msgs, err := ch.Consume(
		queueName,
		"",
		false, // autoAck false
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	go func() {
		for msg := range msgs {
			var reviewJob ReviewJob
			if err := json.Unmarshal(msg.Body, &reviewJob); err != nil {
				log.Printf("json.Unmarshal review job: %v\n", err)
				// If json unmarshaling is not done, do not requeue message.
				if err := msg.Nack(false, false); err != nil {
					log.Printf("msg.Nack: %v\n", err)
				}
				continue
			}
			log.Printf("message received: %s\n", msg.Body)
			if err := msg.Ack(false); err != nil {
				log.Printf("msg.Ack: %v\n", err)
			}
		}
	}()

	return nil
}
