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

type ReviewPublisher interface {
	PublishReviewJob(ctx context.Context, reviewID string, bookID int) error
}

type RabbitMQPublisher struct {
	ch        *amqp.Channel
	queueName string
}

func NewPublisher(ch *amqp.Channel, queueName string) *RabbitMQPublisher {
	return &RabbitMQPublisher{
		ch:        ch,
		queueName: queueName,
	}
}

func (p *RabbitMQPublisher) PublishReviewJob(ctx context.Context, reviewID string, bookID int) error {
	job := ReviewJob{
		ReviewID: reviewID,
		BookID:   bookID,
	}
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	err = p.ch.PublishWithContext(
		ctx,
		"",          // default exchange
		p.queueName, // routing key = queue name
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
