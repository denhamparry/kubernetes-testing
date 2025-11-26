package workload

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
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
		deleteCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := clientset.AppsV1().StatefulSets(namespace).Delete(deleteCtx, statefulSetName, metav1.DeleteOptions{}); err != nil {
			fmt.Printf("Warning: failed to cleanup statefulset %s: %v\n", statefulSetName, err)
		}
	}()

	// Wait for statefulset to be ready using proper wait mechanism
	err = wait.PollUntilContextTimeout(ctx, 1*time.Second, 60*time.Second, true,
		func(ctx context.Context) (bool, error) {
			sts, err := clientset.AppsV1().StatefulSets(namespace).Get(ctx, statefulSetName, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			if sts.Status.ReadyReplicas >= 1 {
				return true, nil
			}
			return false, nil
		})
	if err != nil {
		return fmt.Errorf("statefulset test failed: %w", err)
	}

	return nil
}
