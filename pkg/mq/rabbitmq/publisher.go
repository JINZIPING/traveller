package rabbitmq

import (
	"encoding/json"
	"log"
	"my_project/pkg/mq"

	"github.com/streadway/amqp"
)

// RabbitMQPublisher RabbitMQ的任务发布器实现
type RabbitMQPublisher struct {
	ch         *amqp.Channel
	exchange   string
	queueName  string
	routingKey string
}

// NewRabbitMQPublisher 构造函数
func NewRabbitMQPublisher(ch *amqp.Channel, exchange, queueName, routingKey string) *RabbitMQPublisher {
	return &RabbitMQPublisher{
		ch:         ch,
		exchange:   exchange,
		queueName:  queueName,
		routingKey: routingKey,
	}
}

// Publish 将任务下发到RabbitMQ
func (p *RabbitMQPublisher) Publish(task any) error {
	// 序列化
	body, err := json.Marshal(task)
	if err != nil {
		log.Printf("[RabbitMQPublisher] Marshal task failed: %v", err)
		return err
	}

	// 绑定消息队列
	if err := p.ch.QueueBind(
		p.queueName,
		p.routingKey,
		p.exchange,
		false,
		nil,
	); err != nil {
		log.Printf("[RabbitMQPublisher] QueueBind failed: exchange=%s, queue=%s, key=%s, err=%v",
			p.exchange, p.queueName, p.routingKey, err)
		return err
	}

	//log.Printf("[RabbitMQPublisher] publish task: exchange=%s, routingKey=%s, body=%s",
	//	p.exchange, p.routingKey, string(body))

	// 下发任务
	if err := p.ch.Publish(
		p.exchange,
		p.routingKey,
		false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	); err != nil {
		log.Printf("[RabbitMQPublisher] Publish failed: exchange=%s, key=%s, body=%s, err=%v",
			p.exchange, p.routingKey, string(body), err)
		return err
	}
	log.Printf("[RabbitMQPublisher] Task published successfully: exchange=%s, key=%s, body=%s",
		p.exchange, p.routingKey, string(body))
	return nil
}

// 确认实现接口
var _ mq.TaskPublisher = (*RabbitMQPublisher)(nil)
