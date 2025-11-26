package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ktest",
	Short: "Kubernetes cluster testing tool",
	Long:  `A comprehensive testing tool for Kubernetes clusters including conformance, operational, and performance testing.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("kubeconfig", "", "path to kubeconfig file (default: $HOME/.kube/config)")
}
