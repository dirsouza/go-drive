package queue

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type RabbitMQConfig struct {
	URL       string
	TopicName string
	Timeout   time.Duration
}

type RabbitMQConnection struct {
	cfg  RabbitMQConfig
	conn *amqp.Connection
}

func newRabbitMQConnection(cfg RabbitMQConfig) (rabbit *RabbitMQConnection, err error) {
	rabbit.cfg = cfg
	rabbit.conn, err = amqp.Dial(cfg.URL)
	return
}

func (rabbit *RabbitMQConnection) Publish(msg []byte) error {
	ch, err := rabbit.conn.Channel()
	if err != nil {
		return err
	}

	msgPub := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Timestamp:    time.Now(),
		Body:         msg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return ch.PublishWithContext(ctx, rabbit.cfg.TopicName, "", false, false, msgPub)
}

func (rabbit *RabbitMQConnection) Consume(code chan<- MessageDto) error {
	ch, err := rabbit.conn.Channel()
	if err != nil {
		return err
	}

	queue, err := ch.QueueDeclare(rabbit.cfg.TopicName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	deliveries, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for delivery := range deliveries {
		dto := MessageDto{}
		err = dto.Unmarshal(delivery.Body)
		if err != nil {
			return err
		}
		code <- dto
	}

	return nil
}
