package rmq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ReviewJob struct {
	ReviewID string `json:"review_id"`
	BookID   int    `json:"book_id"`
}

func PublishReviewJob(ch *amqp.Channel, queueName string, job ReviewJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		context.Background(),
		"",        // default exchange
		queueName, // routing key = queue name
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         data,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
