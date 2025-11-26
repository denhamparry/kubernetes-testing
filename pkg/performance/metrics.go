package performance

import (
	"fmt"
	"time"
)

type Metrics struct {
	AverageLatency time.Duration
	P95Latency     time.Duration
	P99Latency     time.Duration
	Throughput     float64
	ErrorRate      float64
	TotalRequests  int
	SuccessCount   int
	FailureCount   int
	Duration       time.Duration
}

func (m *Metrics) Report() string {
	report := fmt.Sprintf(`
Performance Test Results
========================
Duration:         %s
Total Requests:   %d
Successful:       %d
Failed:           %d
Error Rate:       %.2f%%

Latency Metrics:
  Average:        %s
  P95:            %s
  P99:            %s

Throughput:       %.2f req/s
`,
		m.Duration,
		m.TotalRequests,
		m.SuccessCount,
		m.FailureCount,
		m.ErrorRate,
		m.AverageLatency,
		m.P95Latency,
		m.P99Latency,
		m.Throughput,
	)

	return report
}
