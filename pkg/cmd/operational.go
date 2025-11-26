package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/denhamparry/kubernetes-testing/pkg/kubeconfig"
	"github.com/denhamparry/kubernetes-testing/pkg/networking"
	"github.com/denhamparry/kubernetes-testing/pkg/storage"
	"github.com/denhamparry/kubernetes-testing/pkg/workload"
	"github.com/spf13/cobra"
)

var operationalCmd = &cobra.Command{
	Use:   "operational",
	Short: "Run operational tests",
	Long:  `Run operational tests for networking, storage, and workloads`,
	RunE: func(cmd *cobra.Command, args []string) error {
		kubeconfigPath, err := cmd.Flags().GetString("kubeconfig")
		if err != nil {
			return fmt.Errorf("failed to get kubeconfig flag: %w", err)
		}
		tests, err := cmd.Flags().GetStringSlice("tests")
		if err != nil {
			return fmt.Errorf("failed to get tests flag: %w", err)
		}
		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			return fmt.Errorf("failed to get namespace flag: %w", err)
		}

		fmt.Println("Running operational tests...")

		// Load kubeconfig and create client
		client, err := kubeconfig.NewClient(kubeconfigPath)
		if err != nil {
			return fmt.Errorf("failed to create kubernetes client: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		// Determine which tests to run
		runAll := contains(tests, "all")
		runNetworking := runAll || contains(tests, "networking")
		runStorage := runAll || contains(tests, "storage")
		runWorkload := runAll || contains(tests, "workload")

		// Run networking tests
		if runNetworking {
			fmt.Println("\nRunning networking tests...")
			if err := networking.TestDNS(ctx, client.Clientset, namespace); err != nil {
				fmt.Printf("  DNS test: FAILED - %v\n", err)
			} else {
				fmt.Println("  DNS test: PASSED")
			}

			if err := networking.TestPodToPod(ctx, client.Clientset, namespace); err != nil {
				fmt.Printf("  Pod-to-pod connectivity: FAILED - %v\n", err)
			} else {
				fmt.Println("  Pod-to-pod connectivity: PASSED")
			}

			if err := networking.TestServiceConnectivity(ctx, client.Clientset, namespace); err != nil {
				fmt.Printf("  Service connectivity: FAILED - %v\n", err)
			} else {
				fmt.Println("  Service connectivity: PASSED")
			}
		}

		// Run storage tests
		if runStorage {
			fmt.Println("\nRunning storage tests...")
			if err := storage.TestStorageClass(ctx, client.Clientset); err != nil {
				fmt.Printf("  Storage class: FAILED - %v\n", err)
			} else {
				fmt.Println("  Storage class: PASSED")
			}

			if err := storage.TestPVCCreation(ctx, client.Clientset, namespace, ""); err != nil {
				fmt.Printf("  PVC creation: FAILED - %v\n", err)
			} else {
				fmt.Println("  PVC creation: PASSED")
			}
		}

		// Run workload tests
		if runWorkload {
			fmt.Println("\nRunning workload tests...")
			if err := workload.TestDeployment(ctx, client.Clientset, namespace); err != nil {
				fmt.Printf("  Deployment: FAILED - %v\n", err)
			} else {
				fmt.Println("  Deployment: PASSED")
			}

			if err := workload.TestStatefulSet(ctx, client.Clientset, namespace); err != nil {
				fmt.Printf("  StatefulSet: FAILED - %v\n", err)
			} else {
				fmt.Println("  StatefulSet: PASSED")
			}
		}

		fmt.Println("\nOperational tests completed!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(operationalCmd)
	operationalCmd.Flags().StringSlice("tests", []string{"all"}, "Tests to run: networking, storage, workload, all")
	operationalCmd.Flags().String("namespace", "default", "Kubernetes namespace to use for tests")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
