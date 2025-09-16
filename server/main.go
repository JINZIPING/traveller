package main

import (
	"log"
	"my_project/server/config"
	"my_project/server/internal/adapter/metrics"
	"my_project/server/internal/infra"
	"my_project/server/internal/router"
	"my_project/server/internal/service"

	"my_project/pkg/mq/rabbitmq"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/spf13/viper"
)

func main() {
	// 1. 加载配置
	config.InitConfig()

	// 2. 初始化日志
	infra.InitLog()

	// 3. 初始化数据库
	mysqlDB := infra.InitMySQL()
	clickhouseDB := infra.InitClickHouse()
	defer mysqlDB.Close()
	defer clickhouseDB.Close()

	// 4. 初始化 RabbitMQ
	conn, ch, err := infra.InitRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	// 创建 RabbitMQ 工厂
	exchange := viper.GetString("rabbitmq.exchange")

	// ICMP 任务队列 (publisher)
	icmpTaskQueue := viper.GetString("rabbitmq.icmp_task_queue")
	icmpTaskKey := viper.GetString("rabbitmq.icmp_task_routing_key")
	icmpTaskFactory := rabbitmq.NewRabbitMQFactory(ch, exchange, icmpTaskQueue, icmpTaskKey)

	// ICMP 结果队列 (consumer)
	icmpResultQueue := viper.GetString("rabbitmq.icmp_result_queue")
	icmpResultKey := viper.GetString("rabbitmq.icmp_result_routing_key")
	icmpResultFactory := rabbitmq.NewRabbitMQFactory(ch, exchange, icmpResultQueue, icmpResultKey)

	// TCP 任务队列 (publisher)
	tcpTaskQueue := viper.GetString("rabbitmq.tcp_task_queue")
	tcpTaskKey := viper.GetString("rabbitmq.tcp_task_routing_key")
	tcpTaskFactory := rabbitmq.NewRabbitMQFactory(ch, exchange, tcpTaskQueue, tcpTaskKey)

	// TCP 结果队列 (consumer)
	tcpResultQueue := viper.GetString("rabbitmq.tcp_result_queue")
	tcpResultKey := viper.GetString("rabbitmq.tcp_result_routing_key")
	tcpResultFactory := rabbitmq.NewRabbitMQFactory(ch, exchange, tcpResultQueue, tcpResultKey)

	// 4. 初始化 Prometheus Publisher
	promHost := viper.GetString("prometheus.host")
	promPort := viper.GetInt("prometheus.port")
	promJob := viper.GetString("prometheus.job")
	metricsPublisher := metrics.NewPrometheusPublisher(promHost, promPort, promJob)

	// 5. 创建 TaskService
	taskService := service.NewTaskService(
		tcpTaskFactory,
		icmpTaskFactory,
		tcpResultFactory,
		icmpResultFactory,
		metricsPublisher,
	)

	// 6. 初始化 Hertz
	h := server.Default(server.WithHostPorts(":8080"))

	// 7. 注册路由
	router.InitializeRoutes(h, taskService)

	// 消费启动
	go taskService.ConsumeICMPResults()
	go taskService.ConsumeTCPResults()

	// 启动 Hertz 服务器
	log.Println("-----Starting Hertz server on :8080-----")
	h.Spin()
}
