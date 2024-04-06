package prometheus

import (
	"fmt"
	"github.com/falser101/pulsar-operator/api/v1alpha1"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makePrometheusDataVolumeName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-prometheus-data", c.Name)
}

// PV/PVC
func makeVolumeClaimTemplates(c *v1alpha1.PulsarCluster) []v1.PersistentVolumeClaim {
	return []v1.PersistentVolumeClaim{
		makePrometheusDataVolumeClaimTemplate(c),
	}
}

func makePrometheusDataVolumeClaimTemplate(c *v1alpha1.PulsarCluster) v1.PersistentVolumeClaim {
	return v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      makePrometheusDataVolumeName(c),
			Namespace: c.Namespace,
		},
		Spec: makePrometheusDataVolumeClaimSpec(c),
	}
}

func makePrometheusDataVolumeClaimSpec(c *v1alpha1.PulsarCluster) v1.PersistentVolumeClaimSpec {
	capacity := fmt.Sprintf("%dGi", c.Spec.Monitor.Prometheus.StorageCapacity)
	return v1.PersistentVolumeClaimSpec{
		AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
		Resources:        v1.ResourceRequirements{Requests: v1.ResourceList{v1.ResourceStorage: resource.MustParse(capacity)}},
		StorageClassName: &c.Spec.Monitor.Prometheus.StorageClassName,
	}
}

// EmptyDir volume
func makeEmptyDirVolume(c *v1alpha1.PulsarCluster) []v1.Volume {
	return []v1.Volume{
		{
			Name:         makePrometheusDataVolumeName(c),
			VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}},
		},
	}
}
