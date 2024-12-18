package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQConsumer struct {
	channel *amqp.Channel
	queue   string
}

func NewRabbitMQConsumer(conn *amqp.Connection, queueName string) *RabbitMQConsumer {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	return &RabbitMQConsumer{
		channel: ch,
		queue:   q.Name,
	}
}

func (c *RabbitMQConsumer) Consume() (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		c.queue, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
}

func (c *RabbitMQConsumer) Close() {
	c.channel.Close()
}
