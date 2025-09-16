package infra

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

// InitRabbitMQ 初始化 RabbitMQ 连接和通道
func InitRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error

	rabbitmqURL := viper.GetString("rabbitmq.url")
	log.Printf("Initializing RabbitMQ connection with URL: %s", rabbitmqURL)

	// 尝试重连
	for i := 0; i < 10; i++ {
		log.Printf("Attempt %d: Connecting to RabbitMQ...", i+1)

		conn, err = amqp.Dial(rabbitmqURL)
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %v, retrying in 10 seconds...", err)
			time.Sleep(10 * time.Second)
			continue
		}

		log.Println("Successfully connected to RabbitMQ, creating channel...")
		ch, err = conn.Channel()
		if err != nil {
			log.Printf("Failed to create RabbitMQ channel: %v, retrying in 5 seconds...", err)
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("Successfully connected to RabbitMQ and created channel")

		// 声明交换机
		if err := ch.ExchangeDeclare(
			viper.GetString("rabbitmq.exchange"),
			"direct",
			true,  // 持久化
			false, // 不自动删除
			false, // 非排他
			false, // 不阻塞
			nil,
		); err != nil {
			log.Printf("Failed to declare exchange: %v", err)
			ch.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		// 声明队列（你可以按需抽象成数组循环）
		queues := []string{
			viper.GetString("rabbitmq.icmp_task_queue"),
			viper.GetString("rabbitmq.tcp_task_queue"),
			viper.GetString("rabbitmq.icmp_result_queue"),
			viper.GetString("rabbitmq.tcp_result_queue"),
		}

		for _, q := range queues {
			if _, err := ch.QueueDeclare(q, true, false, false, false, nil); err != nil {
				log.Printf("Failed to declare queue %s: %v", q, err)
				ch.Close()
				conn.Close()
				time.Sleep(5 * time.Second)
				continue
			}
			log.Printf("Successfully declared queue: %s", q)
		}

		return conn, ch, nil
	}

	if conn != nil {
		log.Println("Closing RabbitMQ connection due to repeated failures")
		conn.Close()
	}

	return nil, nil, fmt.Errorf("failed to connect to RabbitMQ after multiple attempts: %v", err)
}
