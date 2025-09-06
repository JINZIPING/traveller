package mq

// TaskPublisher
type TaskPublisher interface {
	Publish(task any) error
}
