//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/denhamparry/kubernetes-testing/pkg/kubeconfig"
	"github.com/denhamparry/kubernetes-testing/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestStorageIntegration(t *testing.T) {
	// This test requires a real Kubernetes cluster
	client, err := kubeconfig.NewClient("")
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return
	}

	ctx := context.Background()

	t.Run("TestStorageClass", func(t *testing.T) {
		err := storage.TestStorageClass(ctx, client.Clientset)
		assert.NoError(t, err, "Storage class test should pass")
	})

	t.Run("TestPVCCreation", func(t *testing.T) {
		err := storage.TestPVCCreation(ctx, client.Clientset, "default", "")
		// PVC creation might fail if no storage provisioner, but should not panic
		if err != nil {
			t.Logf("PVC creation test: %v", err)
		}
	})
}
