package service

import (
	"encoding/json"
	"fmt"
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/SmirnovND/gofermart/internal/pkg/rabbitmq"
	"log"
	"time"
)

const (
	exchangeName = "delayed_exchange"
	routingKey   = "delayed_routing_key"
)

type RabbitMqService struct {
	producer *rabbitmq.RabbitMQProducer
	consumer *rabbitmq.RabbitMQConsumer
}

func NewRabbitMqService(
	producer *rabbitmq.RabbitMQProducer,
	consumer *rabbitmq.RabbitMQConsumer,
) *RabbitMqService {
	return &RabbitMqService{
		producer: producer,
		consumer: consumer,
	}
}

func (r *RabbitMqService) SendMessageWithDelay(
	number string,
	userId int,
	delay time.Duration,
) error {
	messageObj := &domain.RabbitMqMessage{
		OrderNumber: number,
		UserId:      userId,
	}
	messageBody, err := json.Marshal(messageObj)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return r.producer.Publish(messageBody, delay, exchangeName, routingKey)
}

func (r *RabbitMqService) Consume(handler func(number string, userId int) error) {
	// Соединение с RabbitMQ и получение сообщений
	msgs, err := r.consumer.Consume()
	if err != nil {
		log.Fatalf("Failed to start consuming messages: %s", err)
	}

	log.Println("Waiting for messages...")
	for d := range msgs {
		fmt.Println("get Consume message")
		// Десериализуем JSON-сообщение в структуру
		var message *domain.RabbitMqMessage
		if err := json.Unmarshal(d.Body, &message); err != nil {
			d.Nack(false, true)
			log.Printf("Failed to unmarshal message: %s", err)
			continue
		}
		fmt.Println("handler Consume message")
		err := handler(message.OrderNumber, message.UserId)
		if err != nil {
			d.Nack(false, true)
			log.Printf("Failed handle message: %s", err)
			continue
		}

		err = d.Ack(false) // Подтверждение успешной обработки
		if err != nil {
			log.Printf("Failed to acknowledge message: %s", err)
		}
	}
}
