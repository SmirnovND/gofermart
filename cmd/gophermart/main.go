package main

import (
	"github.com/SmirnovND/gofermart/internal/config"
	"github.com/SmirnovND/gofermart/internal/container"
	"github.com/SmirnovND/gofermart/internal/middleware"
	"github.com/SmirnovND/gofermart/internal/pkg/compressor"
	"github.com/SmirnovND/gofermart/internal/pkg/logger"
	"github.com/SmirnovND/gofermart/internal/pkg/rabbitmq"
	"github.com/SmirnovND/gofermart/internal/router"
	"github.com/SmirnovND/gofermart/internal/service"
	"github.com/SmirnovND/gofermart/internal/usecase"
	"log"
	"net/http"
)

func main() {
	if err := Run(); err != nil {
		panic(err)
	}
}

func Run() error {
	diContainer := container.NewContainer()

	var cf *config.Config
	diContainer.Invoke(func(c *config.Config) {
		cf = c
	})

	// Подключение к RabbitMQ через контейнер
	var rabbitConnection *rabbitmq.RabbitMQConnection
	diContainer.Invoke(func(rc *rabbitmq.RabbitMQConnection) {
		rabbitConnection = rc
	})
	defer rabbitConnection.Close()

	var (
		rabbitProducer *rabbitmq.RabbitMQProducer
		rabbitConsumer *rabbitmq.RabbitMQConsumer
	)

	err := diContainer.Invoke(func(p *rabbitmq.RabbitMQProducer, c *rabbitmq.RabbitMQConsumer) {
		rabbitProducer = p
		rabbitConsumer = c
	})
	if err != nil {
		log.Fatalf("Failed to resolve dependencies: %s", err)
	}

	defer rabbitProducer.Close()
	defer rabbitConsumer.Close()

	var rabbitMqService *service.RabbitMqService
	var processingUseCase *usecase.ProcessingUseCase
	diContainer.Invoke(func(rs *service.RabbitMqService, pu *usecase.ProcessingUseCase) {
		rabbitMqService = rs
		processingUseCase = pu
	})
	go func() {
		rabbitMqService.Consume(processingUseCase.CheckProcessedAndAccrueBalance)
	}()

	return http.ListenAndServe(cf.GetFlagRunAddr(), middleware.ChainMiddleware(
		router.Handler(diContainer),
		logger.WithLogging,
		compressor.WithDecompression,
		compressor.WithCompression,
	))
}
