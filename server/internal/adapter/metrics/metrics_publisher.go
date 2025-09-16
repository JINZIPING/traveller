package metrics

// MetricsPublisher 指标上报抽象接口
type MetricsPublisher interface {
	PublishMetrics(result any, timestamp int64) error
}
