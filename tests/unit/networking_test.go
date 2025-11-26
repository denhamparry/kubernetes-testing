package unit

import (
	"context"
	"testing"

	"github.com/denhamparry/kubernetes-testing/pkg/networking"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNetworkingFunctions(t *testing.T) {
	// Create a fake clientset for testing
	clientset := fake.NewSimpleClientset()
	ctx := context.Background()

	t.Run("TestDNS", func(t *testing.T) {
		// With fake clientset, the function should execute without panic
		// Testing function signature and basic logic
		err := networking.TestDNS(ctx, clientset, "default")
		// Fake clientset allows pod creation to succeed
		assert.NoError(t, err)
	})

	t.Run("TestPodToPod", func(t *testing.T) {
		err := networking.TestPodToPod(ctx, clientset, "default")
		// Function should execute without panic
		assert.NoError(t, err)
	})

	t.Run("TestServiceConnectivity", func(t *testing.T) {
		err := networking.TestServiceConnectivity(ctx, clientset, "default")
		// Function should execute without panic
		assert.NoError(t, err)
	})
}
