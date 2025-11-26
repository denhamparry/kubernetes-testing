package report

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type TestReport struct {
	TestSuite  string
	StartTime  time.Time
	EndTime    time.Time
	TotalTests int
	Passed     int
	Failed     int
	Skipped    int
	Results    []TestResult
}

type TestResult struct {
	Name     string
	Status   string
	Duration time.Duration
	Message  string
}

func NewTestReport(testSuite string) *TestReport {
	return &TestReport{
		TestSuite: testSuite,
		StartTime: time.Now(),
		Results:   []TestResult{},
	}
}

func (r *TestReport) AddResult(result TestResult) {
	r.Results = append(r.Results, result)
	r.TotalTests++

	switch result.Status {
	case "passed":
		r.Passed++
	case "failed":
		r.Failed++
	case "skipped":
		r.Skipped++
	}
}

func (r *TestReport) Complete() {
	r.EndTime = time.Now()
}

func (r *TestReport) GenerateHTML() (string, error) {
	duration := r.EndTime.Sub(r.StartTime)

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>%s Test Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .summary { background: #f5f5f5; padding: 15px; border-radius: 5px; }
        .passed { color: green; }
        .failed { color: red; }
        .skipped { color: orange; }
        table { width: 100%%; border-collapse: collapse; margin-top: 20px; }
        th, td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #4CAF50; color: white; }
    </style>
</head>
<body>
    <h1>%s Test Report</h1>
    <div class="summary">
        <p><strong>Duration:</strong> %s</p>
        <p><strong>Total Tests:</strong> %d</p>
        <p class="passed"><strong>Passed:</strong> %d</p>
        <p class="failed"><strong>Failed:</strong> %d</p>
        <p class="skipped"><strong>Skipped:</strong> %d</p>
    </div>
    <table>
        <thead>
            <tr>
                <th>Test Name</th>
                <th>Status</th>
                <th>Duration</th>
                <th>Message</th>
            </tr>
        </thead>
        <tbody>
`,
		r.TestSuite, r.TestSuite, duration, r.TotalTests, r.Passed, r.Failed, r.Skipped)

	for _, result := range r.Results {
		statusClass := result.Status
		html += fmt.Sprintf(`
            <tr>
                <td>%s</td>
                <td class="%s">%s</td>
                <td>%s</td>
                <td>%s</td>
            </tr>
`,
			result.Name, statusClass, strings.ToUpper(result.Status), result.Duration, result.Message)
	}

	html += `
        </tbody>
    </table>
</body>
</html>
`

	return html, nil
}

func (r *TestReport) GenerateJSON() (string, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal report to JSON: %w", err)
	}
	return string(data), nil
}

func (r *TestReport) Print() {
	duration := r.EndTime.Sub(r.StartTime)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("%s Test Report\n", r.TestSuite)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Duration:     %s\n", duration)
	fmt.Printf("Total Tests:  %d\n", r.TotalTests)
	fmt.Printf("Passed:       %d\n", r.Passed)
	fmt.Printf("Failed:       %d\n", r.Failed)
	fmt.Printf("Skipped:      %d\n", r.Skipped)
	fmt.Println(strings.Repeat("-", 60))

	for _, result := range r.Results {
		statusSymbol := "✓"
		if result.Status == "failed" {
			statusSymbol = "✗"
		} else if result.Status == "skipped" {
			statusSymbol = "⊘"
		}

		fmt.Printf("%s %s (%s)\n", statusSymbol, result.Name, result.Duration)
		if result.Message != "" {
			fmt.Printf("  %s\n", result.Message)
		}
	}

	fmt.Println(strings.Repeat("=", 60))
}
