package networking

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func TestDNS(ctx context.Context, clientset *kubernetes.Clientset, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	// Create a test pod for DNS resolution
	podName := "dns-test-" + fmt.Sprintf("%d", time.Now().Unix())
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "dns-test",
					Image:   "busybox:latest",
					Command: []string{"sh", "-c", "nslookup kubernetes.default && sleep 3600"},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	// Create the pod
	_, err := clientset.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create DNS test pod: %w", err)
	}

	// Clean up pod on completion
	defer func() {
		clientset.CoreV1().Pods(namespace).Delete(context.Background(), podName, metav1.DeleteOptions{})
	}()

	// Wait for pod to complete (simplified - production would use wait mechanism)
	time.Sleep(10 * time.Second)

	// Check pod status
	testPod, err := clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get DNS test pod status: %w", err)
	}

	if testPod.Status.Phase == corev1.PodFailed {
		return fmt.Errorf("DNS test pod failed")
	}

	return nil
}
