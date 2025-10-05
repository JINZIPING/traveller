package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"my_project/pkg/mark"
	pkgModel "my_project/pkg/model"
	"my_project/pkg/mq"
	"my_project/pkg/utils/timeutil"
	"my_project/server/internal/adapter/metrics"
	"my_project/server/internal/dao"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
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

//// HandleTCPTask 下发TCP探测任务
//func (s *TaskService) HandleTCPTask() app.HandlerFunc {
//	return func(ctx context.Context, c *app.RequestContext) {
//		var req dao.TCPProbeTaskReq
//		if err := c.BindAndValidate(&req); err != nil {
//			c.String(400, "invalid request: %v", err)
//			return
//		}
//
//		if net.ParseIP(req.IP) == nil {
//			c.String(400, "invalid ip")
//			return
//		}
//
//		task := &pkgModel.TCPProbeTaskDTO{
//			IP:        req.IP,
//			Port:      req.Port,
//			Timeout:   req.Timeout,
//			CreatedAt: timeutil.NowUTC8(),
//			UpdatedAt: timeutil.NowUTC8(),
//		}
//
//		if err := s.IssueTCPOnce(ctx, task); err != nil {
//			log.Printf("publish TCP task error: %v", err)
//			c.String(500, "publish error")
//			return
//		}
//
//		// 注册定时任务
//		//if req.IntervalSec > 0 {
//		//	jobID := "tcp:" + req.IP + ":" + req.Port
//		//	s.scheduler.AddTCPJob(jobID, time.Duration(req.IntervalSec)*time.Second, func() {
//		//		_ = s.IssueTCPOnce(ctx, &pkgModel.TCPProbeTaskDTO{
//		//			IP: req.IP, Port: req.Port, Timeout: req.Timeout,
//		//			CreatedAt: timeutil.NowUTC8(), UpdatedAt: timeutil.NowUTC8(),
//		//		})
//		//	})
//		//}
//		c.String(200, "tcp task accepted")
//	}
//}
//
//// HandleICMPTask 下发ICMP探测任务
//func (s *TaskService) HandleICMPTask() app.HandlerFunc {
//	return func(ctx context.Context, c *app.RequestContext) {
//		var req dao.ICMPProbeTaskReq
//		if err := c.BindAndValidate(&req); err != nil {
//			c.String(400, "invalid request: %v", err)
//			return
//		}
//
//		if net.ParseIP(req.IP) == nil {
//			c.String(400, "invalid ip")
//			return
//		}
//
//		task := &pkgModel.ICMPProbeTaskDTO{
//			IP:        req.IP,
//			Count:     req.Count,
//			Threshold: req.Threshold,
//			Timeout:   req.Timeout,
//			CreatedAt: timeutil.NowUTC8(),
//			UpdatedAt: timeutil.NowUTC8(),
//		}
//
//		if err := s.IssueICMPOnce(ctx, task); err != nil {
//			log.Printf("publish ICMP task error: %v", err)
//			c.String(500, "publish error")
//			return
//		}
//
//		//if req.IntervalSec > 0 {
//		//	jobID := "icmp:" + req.IP
//		//	s.scheduler.AddICMPJob(jobID, time.Duration(req.IntervalSec)*time.Second, func() {
//		//		_ = s.IssueICMPOnce(&pkgModel.ICMPProbeTaskDTO{
//		//			IP:        req.IP,
//		//			Count:     req.Count,
//		//			Threshold: req.Threshold,
//		//			Timeout:   req.Timeout,
//		//			CreatedAt: timeutil.NowUTC8(),
//		//			UpdatedAt: timeutil.NowUTC8(),
//		//		})
//		//	})
//		//}
//		c.String(200, "icmp task accepted")
//	}
//}

func (s *TaskService) IssueTCPOnce(ctx context.Context, task *pkgModel.TCPProbeTaskDTO) error {
	if task.TaskID == "" {
		task.TaskID = uuid.NewString()
	}

	// 在 ctx 上触发埋点（记录开始时间）
	ctx = mark.Start(ctx, task.TaskID)

	log.Printf("[MARK] Start TaskID=%s", task.TaskID)

	if err := dao.StoreTCPProbeTask(task); err != nil {
		return err
	}
	pub := s.tcpTaskPublisherFactory.CreatePublisher()
	return pub.Publish(task)
}

func (s *TaskService) IssueICMPOnce(ctx context.Context, task *pkgModel.ICMPProbeTaskDTO) error {
	if task.TaskID == "" {
		task.TaskID = uuid.NewString()
	}
	// 记录开始时间
	ctx = mark.Start(ctx, task.TaskID)

	if err := dao.StoreICMPProbeTask(task); err != nil {
		return err
	}
	pub := s.icmpTaskPublisherFactory.CreatePublisher()
	return pub.Publish(task)
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
		var result pkgModel.TCPProbeResultDTO
		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("unmarshal tcp result failed: %w", err)
		}

		// 用埋点包计算总延时
		totalLatency, ok := mark.Finish(result.TaskID)
		if !ok {
			log.Printf("[WARN] missing start mark for task=%s", result.TaskID)
		}

		// server 收到结果时间，作为指标时间戳
		ts := timeutil.NowUTC8().Unix()

		if err := dao.StoreTCPResult(
			ts,
			result.IP,
			result.Port,
			float64(result.RTT.Microseconds()), result.Success); err != nil {
			return fmt.Errorf("[ERROR] Store tcp result failed: %w", err)
		}
		log.Printf("[TCP MARK] Finish TaskID=%s, latency=%v, ok=%v", result.TaskID, totalLatency, ok)
		return s.metricsPublisher.PublishMetrics(result, ts, totalLatency)
	})
}

// ConsumeICMPResults ICMP消费逻辑
func (s *TaskService) ConsumeICMPResults() {
	s.consumeLoop(s.icmpResultConsumerFactory, func(body []byte) error {
		var result pkgModel.ICMPProbeResultDTO
		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("unmarshal icmp result failed: %w", err)
		}

		// 用埋点包计算总延时
		totalLatency, ok := mark.Finish(result.TaskID)
		if !ok {
			log.Printf("[WARN] missing start mark for task=%s", result.TaskID)
		}

		// server 收到结果时间，作为指标时间戳
		ts := timeutil.NowUTC8().Unix()
		if err := dao.StoreClickHouse(
			ts,
			result.IP,
			result.PacketLoss,
			float64(result.MinRTT.Microseconds()),
			float64(result.MaxRTT.Microseconds()),
			float64(result.AvgRTT.Microseconds()),
		); err != nil {
			return fmt.Errorf("[ERROR] Store icmp result failed: %w", err)
		}
		log.Printf("[ICMP Result] Finish TaskID=%s, latency=%v, ok=%v, Peporting to Prometheus: %f", result.TaskID, totalLatency, ok, result.PacketLoss)
		return s.metricsPublisher.PublishMetrics(result, ts, totalLatency)
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
