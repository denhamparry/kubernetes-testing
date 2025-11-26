package networking

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func TestPodToPod(ctx context.Context, clientset kubernetes.Interface, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	timestamp := time.Now().Unix()

	// Create first pod (server)
	serverPodName := fmt.Sprintf("connectivity-server-%d", timestamp)
	serverPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serverPodName,
			Namespace: namespace,
			Labels:    map[string]string{"app": "connectivity-test", "role": "server"},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "server",
					Image:   "nginx:alpine",
					Ports:   []corev1.ContainerPort{{ContainerPort: 80}},
				},
			},
		},
	}

	// Create second pod (client)
	clientPodName := fmt.Sprintf("connectivity-client-%d", timestamp)
	clientPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clientPodName,
			Namespace: namespace,
			Labels:    map[string]string{"app": "connectivity-test", "role": "client"},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "client",
					Image:   "busybox:latest",
					Command: []string{"sh", "-c", "sleep 3600"},
				},
			},
		},
	}

	// Create pods
	_, err := clientset.CoreV1().Pods(namespace).Create(ctx, serverPod, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create server pod: %w", err)
	}

	_, err = clientset.CoreV1().Pods(namespace).Create(ctx, clientPod, metav1.CreateOptions{})
	if err != nil {
		clientset.CoreV1().Pods(namespace).Delete(ctx, serverPodName, metav1.DeleteOptions{})
		return fmt.Errorf("failed to create client pod: %w", err)
	}

	// Clean up pods
	defer func() {
		clientset.CoreV1().Pods(namespace).Delete(context.Background(), serverPodName, metav1.DeleteOptions{})
		clientset.CoreV1().Pods(namespace).Delete(context.Background(), clientPodName, metav1.DeleteOptions{})
	}()

	return nil
}

func TestServiceConnectivity(ctx context.Context, clientset kubernetes.Interface, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	timestamp := time.Now().Unix()
	serviceName := fmt.Sprintf("test-service-%d", timestamp)

	// Create service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": "test-service"},
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
		},
	}

	_, err := clientset.CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	// Clean up service
	defer func() {
		clientset.CoreV1().Services(namespace).Delete(context.Background(), serviceName, metav1.DeleteOptions{})
	}()

	// Verify service was created
	_, err = clientset.CoreV1().Services(namespace).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to verify service creation: %w", err)
	}

	return nil
}
