package mq

type PublisherFactory interface {
	CreatePublisher() TaskPublisher
}

type ConsumerFactory interface {
	CreateConsumer() ResultConsumer
}
