package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/denhamparry/kubernetes-testing/pkg/performance"
	"github.com/spf13/cobra"
)

var performanceCmd = &cobra.Command{
	Use:   "performance",
	Short: "Run performance tests",
	Long:  `Run performance and load tests against a specified endpoint`,
	RunE: func(cmd *cobra.Command, args []string) error {
		duration, _ := cmd.Flags().GetDuration("duration")
		rps, _ := cmd.Flags().GetInt("rps")
		endpoint, _ := cmd.Flags().GetString("endpoint")

		if endpoint == "" {
			return fmt.Errorf("endpoint is required (use --endpoint flag)")
		}

		fmt.Printf("Running performance test...\n")
		fmt.Printf("  Endpoint: %s\n", endpoint)
		fmt.Printf("  Duration: %s\n", duration)
		fmt.Printf("  Target RPS: %d\n\n", rps)

		// Create and run load test
		loadTest := performance.NewLoadTest(duration, rps, endpoint)
		ctx := context.Background()

		metrics, err := loadTest.Run(ctx)
		if err != nil {
			return fmt.Errorf("performance test failed: %w", err)
		}

		// Print results
		fmt.Println(metrics.Report())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(performanceCmd)
	performanceCmd.Flags().Duration("duration", 5*time.Minute, "Test duration")
	performanceCmd.Flags().Int("rps", 100, "Requests per second")
	performanceCmd.Flags().String("endpoint", "", "Endpoint to test (required)")
}
