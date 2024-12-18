package domain

type RabbitMqMessage struct {
	OrderNumber string `json:"orderNumber"`
	UserId      int    `json:"userId"`
}
