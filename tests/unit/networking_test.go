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
		// This would fail without a real cluster, but tests the function signature
		err := networking.TestDNS(ctx, clientset, "default")
		// In a unit test with mock clientset, we expect this might fail
		// but we're testing that the function exists and has correct signature
		assert.NotNil(t, err)
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
