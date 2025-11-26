package workload

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func TestStatefulSet(ctx context.Context, clientset kubernetes.Interface, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	timestamp := time.Now().Unix()
	statefulSetName := fmt.Sprintf("test-statefulset-%d", timestamp)
	replicas := int32(2)

	// Create statefulset
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulSetName,
			Namespace: namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    &replicas,
			ServiceName: "test-service",
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "test-statefulset",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "test-statefulset",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:alpine",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := clientset.AppsV1().StatefulSets(namespace).Create(ctx, statefulSet, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create statefulset: %w", err)
	}

	// Clean up statefulset
	defer func() {
		clientset.AppsV1().StatefulSets(namespace).Delete(context.Background(), statefulSetName, metav1.DeleteOptions{})
	}()

	// Wait for statefulset to be ready (simplified)
	time.Sleep(10 * time.Second)

	// Verify statefulset
	sts, err := clientset.AppsV1().StatefulSets(namespace).Get(ctx, statefulSetName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get statefulset: %w", err)
	}

	if sts.Status.Replicas == 0 {
		return fmt.Errorf("statefulset has no replicas")
	}

	return nil
}
