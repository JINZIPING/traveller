package rabbitmq

import (
	"my_project/pkg/mq"

	"github.com/streadway/amqp"
)

type RabbitMQFactory struct {
	ch         *amqp.Channel
	exchange   string
	queueName  string
	routingKey string
}

func NewRabbitMQFactory(ch *amqp.Channel, exchange, queueName, routingKey string) *RabbitMQFactory {
	return &RabbitMQFactory{
		ch:         ch,
		exchange:   exchange,
		queueName:  queueName,
		routingKey: routingKey,
	}
}

func (f *RabbitMQFactory) CreatePublisher() mq.TaskPublisher {
	return NewRabbitMQPublisher(f.ch, f.exchange, f.queueName, f.routingKey)
}

func (f *RabbitMQFactory) CreateConsumer() mq.ResultConsumer {
	return NewRabbitMQConsumer(f.ch, f.exchange, f.queueName, f.routingKey)
}
