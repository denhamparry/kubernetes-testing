package storage

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func TestPVCCreation(ctx context.Context, clientset *kubernetes.Clientset, namespace, storageClass string) error {
	if namespace == "" {
		namespace = "default"
	}
	if storageClass == "" {
		storageClass = "standard"
	}

	timestamp := time.Now().Unix()
	pvcName := fmt.Sprintf("test-pvc-%d", timestamp)

	// Create PVC
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcName,
			Namespace: namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
			StorageClassName: &storageClass,
		},
	}

	_, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, pvc, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create PVC: %w", err)
	}

	// Clean up PVC
	defer func() {
		clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(context.Background(), pvcName, metav1.DeleteOptions{})
	}()

	// Wait for PVC to be bound (simplified - production would use proper wait mechanism)
	time.Sleep(5 * time.Second)

	// Verify PVC status
	createdPVC, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get PVC status: %w", err)
	}

	if createdPVC.Status.Phase == corev1.ClaimLost {
		return fmt.Errorf("PVC is in Lost state")
	}

	return nil
}

func TestStorageClass(ctx context.Context, clientset *kubernetes.Clientset) error {
	// List available storage classes
	storageClasses, err := clientset.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list storage classes: %w", err)
	}

	if len(storageClasses.Items) == 0 {
		return fmt.Errorf("no storage classes found in cluster")
	}

	// Check for default storage class
	hasDefault := false
	for _, sc := range storageClasses.Items {
		if sc.Annotations["storageclass.kubernetes.io/is-default-class"] == "true" {
			hasDefault = true
			break
		}
	}

	if !hasDefault {
		return fmt.Errorf("no default storage class found")
	}

	return nil
}
