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
	"net"
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
	scheduler                 *Scheduler
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
		scheduler:                 NewScheduler(),
	}
}

func (s *TaskService) Scheduler() *Scheduler {
	return s.scheduler
}

// HandleTCPTask 下发TCP探测任务
func (s *TaskService) HandleTCPTask() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req dao.TCPProbeTaskReq
		if err := c.BindAndValidate(&req); err != nil {
			c.String(400, "invalid request: %v", err)
			return
		}

		if net.ParseIP(req.IP) == nil {
			c.String(400, "invalid ip")
			return
		}

		task := model.TCPProbeTask{
			IP:        req.IP,
			Port:      req.Port,
			Timeout:   req.Timeout,
			CreatedAt: nowCN(),
			UpdatedAt: nowCN(),
		}

		if err := s.IssueTCPOnce(task); err != nil {
			log.Printf("publish TCP task error: %v", err)
			c.String(500, "publish error")
			return
		}

		// 注册定时任务
		if req.IntervalSec > 0 {
			jobID := "tcp:" + req.IP + ":" + req.Port
			s.scheduler.AddTCPJob(jobID, time.Duration(req.IntervalSec)*time.Second, func() {
				_ = s.IssueTCPOnce(model.TCPProbeTask{
					IP: req.IP, Port: req.Port, Timeout: req.Timeout,
					CreatedAt: nowCN(), UpdatedAt: nowCN(),
				})
			})
		}
		c.String(200, "tcp task accepted")
	}
}

// HandleICMPTask 下发ICMP探测任务
func (s *TaskService) HandleICMPTask() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req dao.ICMPProbeTaskReq
		if err := c.BindAndValidate(&req); err != nil {
			c.String(400, "invalid request: %v", err)
			return
		}

		if net.ParseIP(req.IP) == nil {
			c.String(400, "invalid ip")
			return
		}

		task := model.ICMPProbeTask{
			IP:        req.IP,
			Count:     req.Count,
			Threshold: req.Threshold,
			Timeout:   req.Timeout,
			CreatedAt: nowCN(),
			UpdatedAt: nowCN(),
		}

		if err := s.IssueICMPOnce(task); err != nil {
			log.Printf("publish ICMP task error: %v", err)
			c.String(500, "publish error")
			return
		}

		if req.IntervalSec > 0 {
			jobID := "icmp:" + req.IP
			s.scheduler.AddICMPJob(jobID, time.Duration(req.IntervalSec)*time.Second, func() {
				_ = s.IssueICMPOnce(model.ICMPProbeTask{
					IP:        req.IP,
					Count:     req.Count,
					Threshold: req.Threshold,
					Timeout:   req.Timeout,
					CreatedAt: nowCN(),
					UpdatedAt: nowCN(),
				})
			})
		}
		c.String(200, "icmp task accepted")
	}
}

func (s *TaskService) IssueTCPOnce(task model.TCPProbeTask) error {
	if err := dao.StoreTCPProbeTask(&task); err != nil {
		return err
	}
	pub := s.tcpTaskPublisherFactory.CreatePublisher()
	return pub.Publish(task)
}

func (s *TaskService) IssueICMPOnce(task model.ICMPProbeTask) error {
	if err := dao.StoreICMPProbeTask(&task); err != nil {
		return err
	}
	pub := s.icmpTaskPublisherFactory.CreatePublisher()
	return pub.Publish(task)
}

func nowCN() time.Time {
	loc, _ := time.LoadLocation("Asia/Singapore")
	return time.Now().In(loc)
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
	})
}

// HandlePauseICMP 停止一个 ICMP 周期任务
func (s *TaskService) HandlePauseICMP() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ip := c.Param("ip")
		jobID := "icmp:" + ip
		s.scheduler.Remove(jobID)
		c.String(200, "icmp probe job stopped for %s", ip)
	}
}

// HandlePauseTCP 停止一个 TCP 周期任务
func (s *TaskService) HandlePauseTCP() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ip := c.Param("ip")
		port := c.Param("port")
		jobID := "tcp:" + ip + ":" + port
		s.scheduler.Remove(jobID)
		c.String(200, "tcp probe job stopped for %s:%s", ip, port)
	}
}
