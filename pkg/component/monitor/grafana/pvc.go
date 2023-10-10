package grafana

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cachev1alpha1 "pulsar-operator/pkg/api/v1alpha1"
)

func MakePVC(c *cachev1alpha1.PulsarCluster) *v1.PersistentVolumeClaim {
	capacity := fmt.Sprintf("%dGi", c.Spec.Monitor.Grafana.StorageCapacity)
	return &v1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      makePVCName(c),
			Namespace: c.Namespace,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			Resources:        v1.ResourceRequirements{Requests: v1.ResourceList{v1.ResourceStorage: resource.MustParse(capacity)}},
			StorageClassName: &c.Spec.Monitor.Grafana.StorageClassName,
		},
	}
}

func makePVCName(c *cachev1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-monitor-grafana-data-pvc", c.Name)
}
