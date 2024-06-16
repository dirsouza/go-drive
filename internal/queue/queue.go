package queue

import (
	"fmt"
	"log"
	"reflect"
)

type TypeQueue int

const (
	RabbitMQ TypeQueue = iota
)

type ConnectionQueue interface {
	Publish([]byte) error
	Consume(chan<- MessageDto) error
}

type Queue struct {
	conn ConnectionQueue
}

func New(typeQueue TypeQueue, cfg any) (queue *Queue, err error) {
	rt := reflect.TypeOf(cfg)

	switch typeQueue {
	case RabbitMQ:
		if rt.Name() != "RabbitMQConfig" {
			return nil, fmt.Errorf("config need's to be of type RabbitMQConfig")
		}

		conn, err := newRabbitMQConnection(cfg.(RabbitMQConfig))
		if err != nil {
			return nil, err
		}

		queue.conn = conn
	default:
		log.Fatal("invalid queue type")
	}

	return
}

func (queue *Queue) Publish(msg []byte) error {
	return queue.conn.Publish(msg)
}

func (queue *Queue) Consume(code chan<- MessageDto) error {
	return queue.conn.Consume(code)
}
