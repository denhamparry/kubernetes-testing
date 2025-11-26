package cmd

import (
	"context"
	"fmt"

	"github.com/denhamparry/kubernetes-testing/pkg/conformance"
	"github.com/denhamparry/kubernetes-testing/pkg/kubeconfig"
	"github.com/spf13/cobra"
)

var conformanceCmd = &cobra.Command{
	Use:   "conformance",
	Short: "Run Kubernetes conformance tests",
	Long:  `Run Kubernetes conformance tests using Sonobuoy to validate cluster compliance`,
	RunE: func(cmd *cobra.Command, args []string) error {
		kubeconfigPath, _ := cmd.Flags().GetString("kubeconfig")
		mode, _ := cmd.Flags().GetString("mode")

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
		ctx := context.Background()
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
