package main

import (
	"log"
	"my_project/client/config"
	"my_project/client/internal/infra"
	"my_project/client/internal/service"

	"my_project/pkg/mq/rabbitmq"

	"github.com/spf13/viper"
)

func main() {
	// 1. 加载配置
	config.InitConfig()

	// 2. 初始化日志
	infra.InitLog()

	// 初始化 ICMP RabbitMQ channel
	conn, ch, err := infra.InitRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	exchange := viper.GetString("rabbitmq.exchange")

	// 3. TCP 队列：Consumer (任务) + Publisher (结果)
	tcpTaskQueue := viper.GetString("rabbitmq.tcp_task_queue")
	tcpTaskKey := viper.GetString("rabbitmq.tcp_task_routing_key")
	tcpTaskFactory := rabbitmq.NewRabbitMQFactory(ch, exchange, tcpTaskQueue, tcpTaskKey)

	tcpResultQueue := viper.GetString("rabbitmq.tcp_result_queue")
	tcpResultKey := viper.GetString("rabbitmq.tcp_result_routing_key")
	tcpResultFactory := rabbitmq.NewRabbitMQFactory(ch, exchange, tcpResultQueue, tcpResultKey)

	// 4. 创建 RabbitMQ 工厂（ICMP）
	icmpTaskQueue := viper.GetString("rabbitmq.icmp_task_queue")
	icmpTaskKey := viper.GetString("rabbitmq.icmp_task_routing_key")
	icmpTaskFactory := rabbitmq.NewRabbitMQFactory(ch, exchange, icmpTaskQueue, icmpTaskKey)

	icmpResultQueue := viper.GetString("rabbitmq.icmp_result_queue")
	icmpResultKey := viper.GetString("rabbitmq.icmp_result_routing_key")
	icmpResultFactory := rabbitmq.NewRabbitMQFactory(ch, exchange, icmpResultQueue, icmpResultKey)

	// 5. 初始化 ClientService
	clientService := service.NewClientService(
		tcpTaskFactory,
		tcpResultFactory,
		icmpTaskFactory,
		icmpResultFactory,
	)

	// 6. 启动任务消费
	go clientService.ConsumeTCPTasks()
	go clientService.ConsumeICMPTasks()

	// 7. 阻塞 main，不退出
	select {}
}
