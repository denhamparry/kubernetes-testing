package unit

import (
	"testing"
	"time"

	"github.com/denhamparry/kubernetes-testing/pkg/report"
	"github.com/stretchr/testify/assert"
)

func TestReportGeneration(t *testing.T) {
	t.Run("NewTestReport", func(t *testing.T) {
		r := report.NewTestReport("Conformance")
		assert.NotNil(t, r)
		assert.Equal(t, "Conformance", r.TestSuite)
		assert.Equal(t, 0, r.TotalTests)
	})

	t.Run("AddResult", func(t *testing.T) {
		r := report.NewTestReport("Operational")

		r.AddResult(report.TestResult{
			Name:     "DNS Test",
			Status:   "passed",
			Duration: 5 * time.Second,
			Message:  "All checks passed",
		})

		assert.Equal(t, 1, r.TotalTests)
		assert.Equal(t, 1, r.Passed)
		assert.Equal(t, 0, r.Failed)
	})

	t.Run("GenerateJSON", func(t *testing.T) {
		r := report.NewTestReport("Performance")
		r.AddResult(report.TestResult{
			Name:     "Load Test",
			Status:   "passed",
			Duration: 1 * time.Minute,
			Message:  "Performance within acceptable limits",
		})
		r.Complete()

		jsonStr, err := r.GenerateJSON()
		assert.NoError(t, err)
		assert.Contains(t, jsonStr, "Performance")
		assert.Contains(t, jsonStr, "Load Test")
	})

	t.Run("GenerateHTML", func(t *testing.T) {
		r := report.NewTestReport("Integration")
		r.AddResult(report.TestResult{
			Name:     "API Test",
			Status:   "failed",
			Duration: 2 * time.Second,
			Message:  "Connection timeout",
		})
		r.Complete()

		html, err := r.GenerateHTML()
		assert.NoError(t, err)
		assert.Contains(t, html, "<!DOCTYPE html>")
		assert.Contains(t, html, "Integration Test Report")
		assert.Contains(t, html, "API Test")
	})
}
