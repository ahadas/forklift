package utils

import (
	"context"
	"fmt"
	"github.com/onsi/ginkgo"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	// DefaultStorageClass the default storage class used in tests
	DefaultStorageClass *storagev1.StorageClass
)

// ClientsIface is the clients interface
type ClientsIface interface {
	K8s() *kubernetes.Clientset
}

// GetDefaultStorageClass return the storage class which is marked as default in the cluster
func GetDefaultStorageClass(client *kubernetes.Clientset) *storagev1.StorageClass {
	storageclasses, err := client.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		ginkgo.Fail("Unable to list storage classes")
		return nil
	}
	for _, storageClass := range storageclasses.Items {
		if storageClass.Annotations["storageclass.kubernetes.io/is-default-class"] == "true" {
			return &storageClass
		}
	}
	ginkgo.Fail("Unable to find default storage classes")
	return nil
}

func getDefaultStorageClassCsiDriver(client *kubernetes.Clientset) *storagev1.CSIDriver {
	if DefaultStorageClass != nil {
		csidrivers, err := client.StorageV1().CSIDrivers().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			ginkgo.Fail(fmt.Sprintf("Unable to get csi driver: %v", err))
		}
		for _, driver := range csidrivers.Items {
			if driver.Name == DefaultStorageClass.Provisioner {
				return &driver
			}
		}
	}
	return nil
}

// CacheTestsData fetch and cache data required for tests
func CacheTestsData(client *kubernetes.Clientset, cdiNs string) {
	if DefaultStorageClass == nil {
		DefaultStorageClass = GetDefaultStorageClass(client)
	}
}
