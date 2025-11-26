package unit

import (
	"context"
	"testing"

	"github.com/denhamparry/kubernetes-testing/pkg/networking"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestNetworkingFunctions(t *testing.T) {
	// Create a fake clientset for testing
	clientset := fake.NewSimpleClientset()

	// Add a reactor to simulate pod status changes
	clientset.PrependReactor("get", "pods", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		getAction := action.(k8stesting.GetAction)
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      getAction.GetName(),
				Namespace: getAction.GetNamespace(),
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodSucceeded,
			},
		}
		return true, pod, nil
	})

	ctx := context.Background()

	t.Run("TestDNS", func(t *testing.T) {
		// With fake clientset and reactor, the function should complete successfully
		err := networking.TestDNS(ctx, clientset, "default")
		// Fake clientset allows pod creation and status checking to succeed
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
