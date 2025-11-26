package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/denhamparry/kubernetes-testing/pkg/kubeconfig"
	"github.com/stretchr/testify/assert"
)

func TestNewClient_InvalidPath(t *testing.T) {
	// Test with non-existent kubeconfig file
	_, err := kubeconfig.NewClient("/nonexistent/path/to/kubeconfig")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to build config from kubeconfig")
}

func TestNewClient_EmptyPath(t *testing.T) {
	// Test with empty path - should try default location
	// This test will fail if no kubeconfig exists, which is expected in CI
	_, err := kubeconfig.NewClient("")
	if err != nil {
		// Expected in environments without kubeconfig
		home, _ := os.UserHomeDir()
		defaultPath := filepath.Join(home, ".kube", "config")
		assert.Contains(t, err.Error(), "failed to build config from kubeconfig",
			"Error should indicate kubeconfig issue for path: %s", defaultPath)
	}
}
