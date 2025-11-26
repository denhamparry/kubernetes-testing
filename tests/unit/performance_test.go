package unit

import (
	"testing"
	"time"

	"github.com/denhamparry/kubernetes-testing/pkg/performance"
	"github.com/stretchr/testify/assert"
)

func TestLoadTestMetrics(t *testing.T) {
	t.Run("NewLoadTest", func(t *testing.T) {
		lt := performance.NewLoadTest(5*time.Minute, 100, "http://example.com")
		assert.NotNil(t, lt)
		assert.Equal(t, 5*time.Minute, lt.Duration)
		assert.Equal(t, 100, lt.RPS)
		assert.Equal(t, "http://example.com", lt.Endpoint)
	})

	t.Run("MetricsReport", func(t *testing.T) {
		metrics := &performance.Metrics{
			AverageLatency: 100 * time.Millisecond,
			P95Latency:     200 * time.Millisecond,
			P99Latency:     300 * time.Millisecond,
			Throughput:     50.5,
			ErrorRate:      2.5,
			TotalRequests:  1000,
			SuccessCount:   975,
			FailureCount:   25,
			Duration:       1 * time.Minute,
		}

		report := metrics.Report()
		assert.Contains(t, report, "Performance Test Results")
		assert.Contains(t, report, "1000")
		assert.Contains(t, report, "975")
		assert.Contains(t, report, "25")
	})
}
