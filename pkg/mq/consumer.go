package mq

// ResultConsumer 结果消费者接口
type ResultConsumer interface {
	Consume() (<-chan []byte, error)
}
