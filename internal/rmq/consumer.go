package rmq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ReviewJobProcessor interface {
	EnrichReview(ctx context.Context, reviewID string, bookID int) error
}

type RabbitMQConsumer struct {
	ch        *amqp.Channel
	queueName string
	processor ReviewJobProcessor
}

func NewConsumer(ch *amqp.Channel, queueName string, processor ReviewJobProcessor) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		ch:        ch,
		queueName: queueName,
		processor: processor,
	}
}

func (r *RabbitMQConsumer) StartReviewConsumer() error {
	msgs, err := r.ch.Consume(
		r.queueName,
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

			// sleeping 10 seconds to check review "pending" status before consuming the rabbit queue.
			time.Sleep(10 * time.Second)
			err := r.processor.EnrichReview(context.Background(), reviewJob.ReviewID, reviewJob.BookID)
			if err != nil {
				if err := msg.Nack(false, false); err != nil {
					log.Printf("msg.Nack: %v\n", err)
				}
				continue
			}

			if err := msg.Ack(false); err != nil {
				log.Printf("msg.Ack: %v\n", err)
			}
		}
	}()

	return nil
}
