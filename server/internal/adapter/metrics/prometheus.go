package metrics

import (
	"bytes"
	"fmt"
	"log"
	pkgModel "my_project/pkg/model"
	"net/http"
	"time"
)

type PrometheusPublisher struct {
	host string
	port int
	job  string
}

// NewPrometheusPublisher 构造函数
func NewPrometheusPublisher(host string, port int, job string) *PrometheusPublisher {
	return &PrometheusPublisher{
		host: host,
		port: port,
		job:  job,
	}
}

func (p *PrometheusPublisher) PublishMetrics(result any, timestamp int64, totalLatency time.Duration) error {
	uri := fmt.Sprintf("http://%s:%d/metrics/job/%s", p.host, p.port, p.job)
	var metrics string
	switch r := result.(type) {
	case pkgModel.TCPProbeResultDTO:
		metrics = fmt.Sprintf(
			"tcp_rtt{ip=\"%s\", port=\"%s\"} %f\n"+"tcp_success{ip=\"%s\", port=\"%s\"} %d\n"+
				"tcp_last_seen{ip=\"%s\", port=\"%s\"} %d\n"+"tcp_latency{ip=\"%s\", port=\"%s\"} %f\n",
			r.IP, r.Port, r.RTT.Seconds()*1000, r.IP, r.Port, boolToInt(r.Success),
			r.IP, r.Port, timestamp, r.IP, r.Port, totalLatency.Seconds()*1000,
		)
	case pkgModel.ICMPProbeResultDTO: // ICMP
		metrics = fmt.Sprintf(
			"icmp_packet_loss{ip=\"%s\"} %f\n"+"icmp_rtt_min{ip=\"%s\"} %f\n"+"icmp_rtt_max{ip=\"%s\"} %f\n"+
				"icmp_rtt_avg{ip=\"%s\"} %f\n"+"icmp_last_seen{ip=\"%s\"} %d\n"+"icmp_latency{ip=\"%s\"} %f\n",
			r.IP, r.PacketLoss,
			r.IP, float64(r.MinRTT.Microseconds())/1000.0, r.IP, float64(r.MaxRTT.Microseconds())/1000.0,
			r.IP, float64(r.AvgRTT.Microseconds())/1000.0, r.IP, timestamp, r.IP, totalLatency.Seconds()*1000,
		)
	default:
		return fmt.Errorf("unsupported result type: %T", r)
	}

	log.Printf("[PrometheusPublisher] send metrics:\n%s", metrics)

	req, err := http.NewRequest("POST", uri, bytes.NewBufferString(metrics))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("failed to push metrics, status: %s", resp.Status)
	}

	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
