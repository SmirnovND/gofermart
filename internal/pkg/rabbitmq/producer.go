package rabbitmq

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQProducer struct {
	channel *amqp.Channel
	queue   string
}

func NewRabbitMQProducer(conn *amqp.Connection) *RabbitMQProducer {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	err = ch.ExchangeDeclare(
		"delayed_exchange",
		"x-delayed-message",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-delayed-type": "direct",
		},
	)
	if err != nil {
		return nil
	}

	return &RabbitMQProducer{
		channel: ch,
	}
}

func (p *RabbitMQProducer) Publish(message []byte, delay time.Duration, exchange string, key string) error {
	headers := amqp.Table{
		"x-delay": int32(delay.Milliseconds()),
	}

	return p.channel.Publish(
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			Body:    message,
			Headers: headers,
		},
	)
}

func (p *RabbitMQProducer) Close() {
	p.channel.Close()
}
