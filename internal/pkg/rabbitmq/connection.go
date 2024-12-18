package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

// RabbitMQConnection - общая структура для подключения
type RabbitMQConnection struct {
	Conn *amqp.Connection
}

// NewRabbitMQConnection - создание общего подключения к RabbitMQ
func NewRabbitMQConnection(url string) *RabbitMQConnection {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	return &RabbitMQConnection{Conn: conn}
}

// Close - закрытие подключения
func (r *RabbitMQConnection) Close() {
	if r.Conn != nil {
		r.Conn.Close()
	}
}
