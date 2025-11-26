package storage

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

func TestPVCCreation(ctx context.Context, clientset kubernetes.Interface, namespace, storageClass string) error {
	if namespace == "" {
		namespace = "default"
	}
	if storageClass == "" {
		// Auto-detect default storage class instead of hardcoding "standard"
		defaultSC, err := getDefaultStorageClass(ctx, clientset)
		if err != nil {
			return fmt.Errorf("failed to detect default storage class: %w", err)
		}
		storageClass = defaultSC
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
		deleteCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(deleteCtx, pvcName, metav1.DeleteOptions{}); err != nil {
			fmt.Printf("Warning: failed to cleanup PVC %s: %v\n", pvcName, err)
		}
	}()

	// Wait for PVC to be bound using proper wait mechanism
	err = wait.PollUntilContextTimeout(ctx, 1*time.Second, 60*time.Second, true,
		func(ctx context.Context) (bool, error) {
			pvc, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			if pvc.Status.Phase == corev1.ClaimBound {
				return true, nil
			}
			if pvc.Status.Phase == corev1.ClaimLost {
				return false, fmt.Errorf("PVC is in Lost state")
			}
			return false, nil
		})
	if err != nil {
		return fmt.Errorf("PVC test failed: %w", err)
	}

	return nil
}

func TestStorageClass(ctx context.Context, clientset kubernetes.Interface) error {
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

// getDefaultStorageClass finds and returns the name of the default storage class in the cluster.
func getDefaultStorageClass(ctx context.Context, clientset kubernetes.Interface) (string, error) {
	storageClasses, err := clientset.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list storage classes: %w", err)
	}

	if len(storageClasses.Items) == 0 {
		return "", fmt.Errorf("no storage classes found in cluster")
	}

	// Look for default storage class
	for _, sc := range storageClasses.Items {
		if sc.Annotations["storageclass.kubernetes.io/is-default-class"] == "true" {
			return sc.Name, nil
		}
	}

	// If no default found, return error instead of assuming
	return "", fmt.Errorf("no default storage class found (consider setting --storage-class flag)")
}
