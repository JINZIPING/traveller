package rabbitmq

import (
	"log"
	"my_project/pkg/mq"

	"github.com/streadway/amqp"
)

// RabbitMQConsumer RabbitMQ结果处理器实现
type RabbitMQConsumer struct {
	ch         *amqp.Channel
	exchange   string
	queueName  string
	routingKey string
}

func NewRabbitMQConsumer(ch *amqp.Channel, exchange, queueName, routingKey string) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		ch:         ch,
		exchange:   exchange,
		queueName:  queueName,
		routingKey: routingKey,
	}
}

func (c *RabbitMQConsumer) Consume() (<-chan []byte, error) {
	// 确保绑定
	if err := c.ch.QueueBind(
		c.queueName,
		c.routingKey,
		c.exchange,
		false,
		nil,
	); err != nil {
		log.Printf("[RabbitMQConsumer] QueueBind failed: exchange=%s, queue=%s, key=%s, err=%v",
			c.exchange, c.queueName, c.routingKey, err)
		return nil, err
	}

	msgs, err := c.ch.Consume(
		c.queueName,
		"",
		true,  // auto-ack
		false, // not exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Printf("[RabbitMQConsumer] Consume failed: queue=%s, err=%v", c.queueName, err)
		return nil, err
	}

	out := make(chan []byte)
	go func() {
		defer close(out)
		for msg := range msgs {
			out <- msg.Body
		}
	}()
	return out, nil
}

// 确认实现接口
var _ mq.ResultConsumer = (*RabbitMQConsumer)(nil)
