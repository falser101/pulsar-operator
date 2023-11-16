package grafana

import (
	"fmt"
	"github.com/falser101/pulsar-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakePVC(c *v1alpha1.Pulsar) *v1.PersistentVolumeClaim {
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

func makePVCName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-monitor-grafana-data-pvc", c.Name)
}
