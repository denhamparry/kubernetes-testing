package performance

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"
)

type LoadTest struct {
	Duration time.Duration
	RPS      int
	Endpoint string
}

func NewLoadTest(duration time.Duration, rps int, endpoint string) *LoadTest {
	return &LoadTest{
		Duration: duration,
		RPS:      rps,
		Endpoint: endpoint,
	}
}

func (l *LoadTest) Run(ctx context.Context) (*Metrics, error) {
	if l.Endpoint == "" {
		return nil, fmt.Errorf("endpoint cannot be empty")
	}
	if l.RPS <= 0 {
		return nil, fmt.Errorf("RPS must be greater than 0")
	}
	if l.Duration <= 0 {
		return nil, fmt.Errorf("duration must be greater than 0")
	}

	startTime := time.Now()
	var totalRequests int
	var successfulRequests int
	var failedRequests int
	var latencies []time.Duration

	// Rate limiter
	ticker := time.NewTicker(time.Second / time.Duration(l.RPS))
	defer ticker.Stop()

	timeout := time.After(l.Duration)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for {
		select {
		case <-ctx.Done():
			return l.calculateMetrics(startTime, totalRequests, successfulRequests, failedRequests, latencies), nil
		case <-timeout:
			return l.calculateMetrics(startTime, totalRequests, successfulRequests, failedRequests, latencies), nil
		case <-ticker.C:
			reqStart := time.Now()
			resp, err := client.Get(l.Endpoint)
			reqDuration := time.Since(reqStart)

			totalRequests++
			latencies = append(latencies, reqDuration)

			if err != nil {
				failedRequests++
			} else {
				resp.Body.Close()
				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					successfulRequests++
				} else {
					failedRequests++
				}
			}
		}
	}
}

func (l *LoadTest) calculateMetrics(startTime time.Time, total, success, failed int, latencies []time.Duration) *Metrics {
	duration := time.Since(startTime)
	avgLatency, p95, p99 := calculateLatencyPercentiles(latencies)

	return &Metrics{
		AverageLatency: avgLatency,
		P95Latency:     p95,
		P99Latency:     p99,
		Throughput:     float64(total) / duration.Seconds(),
		ErrorRate:      float64(failed) / float64(total) * 100,
		TotalRequests:  total,
		SuccessCount:   success,
		FailureCount:   failed,
		Duration:       duration,
	}
}

func calculateLatencyPercentiles(latencies []time.Duration) (avg, p95, p99 time.Duration) {
	if len(latencies) == 0 {
		return 0, 0, 0
	}

	// Sort latencies for correct percentile calculation
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	// Calculate average
	var sum time.Duration
	for _, l := range latencies {
		sum += l
	}
	avg = sum / time.Duration(len(latencies))

	// Calculate percentiles from sorted array
	p95Index := int(float64(len(latencies)) * 0.95)
	p99Index := int(float64(len(latencies)) * 0.99)

	if p95Index >= len(latencies) {
		p95Index = len(latencies) - 1
	}
	if p99Index >= len(latencies) {
		p99Index = len(latencies) - 1
	}

	p95 = latencies[p95Index]
	p99 = latencies[p99Index]

	return avg, p95, p99
}
