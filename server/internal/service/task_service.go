package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"my_project/pkg/model"
	"my_project/pkg/mq"
	"my_project/server/internal/adapter/metrics"
	"my_project/server/internal/dao"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

// TaskService 负责任务下发 结果处理 结果上报
type TaskService struct {
	tcpTaskPublisherFactory   mq.PublisherFactory
	icmpTaskPublisherFactory  mq.PublisherFactory
	tcpResultConsumerFactory  mq.ConsumerFactory
	icmpResultConsumerFactory mq.ConsumerFactory
	metricsPublisher          metrics.MetricsPublisher
}

// NewTaskService 构造函数
func NewTaskService(
	tcpPubF mq.PublisherFactory,
	icmpPubF mq.PublisherFactory,
	tcpConF mq.ConsumerFactory,
	icmpConF mq.ConsumerFactory,
	mp metrics.MetricsPublisher,
) *TaskService {
	return &TaskService{
		tcpTaskPublisherFactory:   tcpPubF,
		icmpTaskPublisherFactory:  icmpPubF,
		tcpResultConsumerFactory:  tcpConF,
		icmpResultConsumerFactory: icmpConF,
		metricsPublisher:          mp,
	}
}

// HandleTCPTask 下发TCP探测任务
func (s *TaskService) HandleTCPTask() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		task := model.TCPProbeTask{
			IP:        "1.1.1.1",
			Port:      "80",
			Timeout:   5,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := dao.StoreTCPProbeTask(&task)
		if err != nil {
			log.Printf("Error storing tcp probe task: %v", err)
			c.String(500, "Failed to insert TCP task into MySQL")
			return
		}

		pub := s.tcpTaskPublisherFactory.CreatePublisher()
		if err := pub.Publish(task); err != nil {
			log.Printf("Failed to publish TCP task to queue: %v", err)
			c.String(500, "publish error: %v", err)
			return
		}
		log.Println("TCP task published successfully to RabbitMQ")
		c.String(200, "task assigned successfully")
	}
}

// HandleICMPTask 下发ICMP探测任务
func (s *TaskService) HandleICMPTask() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		task := model.ICMPProbeTask{
			IP:        "1.1.1.1", // 示例IP，实际可从请求或数据库中获取
			Count:     4,
			Threshold: 10, // 丢包率阈值
			Timeout:   5,  // 超时时间
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := dao.StoreICMPProbeTask(&task)
		if err != nil {
			log.Printf("Error storing tcp probe task: %v", err)
			c.String(500, "Failed to insert ICMP task into MySQL")
			return
		}

		pub := s.icmpTaskPublisherFactory.CreatePublisher()
		if err := pub.Publish(task); err != nil {
			log.Printf("Failed to publish ICMP task to queue: %v", err)
			c.String(500, "publish error: %v", err)
			return
		}
		c.String(200, "task assigned successfully")
		log.Println("ICMP task published successfully to RabbitMQ")
	}
}

// consumeLoop 通用处理消息循环
func (s *TaskService) consumeLoop(factory mq.ConsumerFactory, processor func([]byte) error) {
	consumer := factory.CreateConsumer()
	msgs, err := consumer.Consume()
	if err != nil {
		log.Fatalf("Failed to consume: %v", err)
	}
	go func() {
		for msg := range msgs {
			if err := processor(msg); err != nil {
				log.Printf("Failed to process message: %v", err)
			}
		}
	}()
}

// ConsumeTCPResults TCP消费逻辑
func (s *TaskService) ConsumeTCPResults() {
	s.consumeLoop(s.tcpResultConsumerFactory, func(body []byte) error {
		var result model.TCPProbeResult
		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("unmarshal tcp result failed: %w", err)
		}
		log.Printf("Processed TCP probe result: %+v", result)

		ts := result.Timestamp.Unix()
		if err := dao.StoreTCPResult(
			ts,
			result.IP,
			result.Port,
			float64(result.RTT.Microseconds()), result.Success); err != nil {
			return fmt.Errorf("store tcp result failed: %w", err)
		}
		return s.metricsPublisher.PublishMetrics(result, ts)
	})
}

// ConsumeICMPResults ICMP消费逻辑
func (s *TaskService) ConsumeICMPResults() {
	s.consumeLoop(s.icmpResultConsumerFactory, func(body []byte) error {
		var result model.ICMPProbeResult
		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("unmarshal icmp result failed: %w", err)
		}
		log.Printf("Processed ICMP probe result: %+v", result)

		ts := result.Timestamp.Unix()
		if err := dao.StoreClickHouse(
			ts,
			result.IP,
			result.PacketLoss,
			float64(result.MinRTT.Microseconds()),
			float64(result.MaxRTT.Microseconds()),
			float64(result.AvgRTT.Microseconds()),
		); err != nil {
			return fmt.Errorf("store icmp result failed: %w", err)
		}
		log.Printf("Peporting to Prometheus: %f", result.PacketLoss)
		return s.metricsPublisher.PublishMetrics(result, ts)
		//if result.PacketLoss > float64(result.Threshold) {
		//	log.Printf("Packet loss exceeds threshold, reporting to Prometheus: %f > %f", result.PacketLoss, result.Threshold)
		//	return s.metricsPublisher.PublishMetrics(result, ts)
		//}
		//return nil
	})
}
