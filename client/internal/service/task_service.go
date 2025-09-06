package service

import (
	"encoding/json"
	"log"
	"my_project/client/internal/probe"
	"my_project/pkg/model"
	"my_project/pkg/mq"
)

type ClientService struct {
	tcpTaskFactory    mq.ConsumerFactory
	tcpResultFactory  mq.PublisherFactory
	icmpTaskFactory   mq.ConsumerFactory
	icmpResultFactory mq.PublisherFactory
}

func NewClientService(
	tcpTaskF mq.ConsumerFactory,
	tcpResultF mq.PublisherFactory,
	icmpTaskF mq.ConsumerFactory,
	icmpResultF mq.PublisherFactory,
) *ClientService {
	return &ClientService{
		tcpTaskFactory:    tcpTaskF,
		tcpResultFactory:  tcpResultF,
		icmpTaskFactory:   icmpTaskF,
		icmpResultFactory: icmpResultF,
	}
}

// 通用消费循环
func (s *ClientService) consumeLoop(factory mq.ConsumerFactory, processor func([]byte) error) {
	consumer := factory.CreateConsumer()
	msgs, err := consumer.Consume()
	if err != nil {
		log.Fatalf("Failed to consume: %v", err)
	}
	// 协程
	go func() {
		for msg := range msgs {
			if err := processor(msg); err != nil {
				log.Printf("Failed to process task: %v", err)
			}
		}
	}()
}

// ConsumeTCPTasks 任务消费逻辑
func (s *ClientService) ConsumeTCPTasks() {
	s.consumeLoop(s.tcpTaskFactory, func(body []byte) error {
		var task model.TCPProbeTask
		if err := json.Unmarshal(body, &task); err != nil {
			log.Printf("Failed to unmarshal TCP task: %v", err)
			return err
		}

		// 执行探测
		result := probe.ExecuteTCPProbeTask(&task)

		// 上报结果
		pub := s.tcpResultFactory.CreatePublisher()
		if err := pub.Publish(result); err != nil {
			log.Printf("Failed to publish TCP result: %v", err)
			return err
		}
		return nil
	})
}

// ConsumeICMPTasks 任务消费逻辑
func (s *ClientService) ConsumeICMPTasks() {
	s.consumeLoop(s.icmpTaskFactory, func(body []byte) error {
		var task model.ICMPProbeTask
		if err := json.Unmarshal(body, &task); err != nil {
			log.Printf("Failed to unmarshal ICMP task: %v", err)
			return err
		}

		// 执行探测
		result := probe.ExecuteICMPProbeTask(&task)

		// 上报结果
		pub := s.icmpResultFactory.CreatePublisher()
		if err := pub.Publish(result); err != nil {
			log.Printf("Failed to publish ICMP result: %v", err)
			return err
		}
		return nil
	})
}
