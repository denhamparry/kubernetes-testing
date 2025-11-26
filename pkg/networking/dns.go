package networking

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

func TestDNS(ctx context.Context, clientset kubernetes.Interface, namespace string) error {
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
		deleteCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := clientset.CoreV1().Pods(namespace).Delete(deleteCtx, podName, metav1.DeleteOptions{}); err != nil {
			fmt.Printf("Warning: failed to cleanup pod %s: %v\n", podName, err)
		}
	}()

	// Wait for pod to complete using proper wait mechanism
	err = wait.PollUntilContextTimeout(ctx, 1*time.Second, 60*time.Second, true,
		func(ctx context.Context) (bool, error) {
			pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			if pod.Status.Phase == corev1.PodSucceeded {
				return true, nil
			}
			if pod.Status.Phase == corev1.PodFailed {
				return false, fmt.Errorf("DNS test pod failed")
			}
			return false, nil
		})
	if err != nil {
		return fmt.Errorf("DNS test failed: %w", err)
	}

	return nil
}
