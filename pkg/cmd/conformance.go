package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/denhamparry/kubernetes-testing/pkg/conformance"
	"github.com/denhamparry/kubernetes-testing/pkg/kubeconfig"
	"github.com/spf13/cobra"
)

var conformanceCmd = &cobra.Command{
	Use:   "conformance",
	Short: "Run Kubernetes conformance tests",
	Long:  `Run Kubernetes conformance tests using Sonobuoy to validate cluster compliance`,
	RunE: func(cmd *cobra.Command, args []string) error {
		kubeconfigPath, err := cmd.Flags().GetString("kubeconfig")
		if err != nil {
			return fmt.Errorf("failed to get kubeconfig flag: %w", err)
		}
		mode, err := cmd.Flags().GetString("mode")
		if err != nil {
			return fmt.Errorf("failed to get mode flag: %w", err)
		}

		fmt.Printf("Running conformance tests in %s mode...\n", mode)

		// Load kubeconfig and create client
		client, err := kubeconfig.NewClient(kubeconfigPath)
		if err != nil {
			return fmt.Errorf("failed to create kubernetes client: %w", err)
		}

		// Create conformance test
		test, err := conformance.NewConformanceTest(client.Clientset)
		if err != nil {
			return fmt.Errorf("failed to create conformance test: %w", err)
		}

		// Run conformance test
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()
		if err := test.Run(ctx, mode); err != nil {
			return fmt.Errorf("conformance test failed: %w", err)
		}

		// Get results
		results, err := test.GetResults()
		if err != nil {
			return fmt.Errorf("failed to get test results: %w", err)
		}

		fmt.Printf("Conformance test completed:\n")
		fmt.Printf("  Status: %s\n", results.Status)
		fmt.Printf("  Passed: %d\n", results.Passed)
		fmt.Printf("  Failed: %d\n", results.Failed)
		fmt.Printf("  Duration: %s\n", results.Duration)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(conformanceCmd)
	conformanceCmd.Flags().String("mode", "quick", "Test mode: quick or certified-conformance")
}
