package zookeeper

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cachev1alpha1 "pulsar-operator/pkg/api/v1alpha1"
)

// PV/PVC
func makeVolumeClaimTemplates(c *cachev1alpha1.PulsarCluster) []v1.PersistentVolumeClaim {
	return []v1.PersistentVolumeClaim{
		makeDataVolumeClaimTemplate(c),
	}
}

func makeDataVolumeClaimTemplate(c *cachev1alpha1.PulsarCluster) v1.PersistentVolumeClaim {
	return v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      makeDataName(c),
			Namespace: c.Namespace,
		},
		Spec: makeDataVolumeClaimSpec(c),
	}
}

func makeDataName(c *cachev1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-zookeeper-pvc", c.Name)
}

func makeDataVolumeClaimSpec(c *cachev1alpha1.PulsarCluster) v1.PersistentVolumeClaimSpec {
	capacity := fmt.Sprintf("%dGi", c.Spec.Zookeeper.StorageCapacity)
	return v1.PersistentVolumeClaimSpec{
		AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
		Resources:        v1.ResourceRequirements{Requests: v1.ResourceList{v1.ResourceStorage: resource.MustParse(capacity)}},
		StorageClassName: &c.Spec.Zookeeper.StorageClassName,
	}
}

func makeEmptyDirVolume(c *cachev1alpha1.PulsarCluster) []v1.Volume {
	return []v1.Volume{
		{
			Name:         ContainerDataVolumeName,
			VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}},
		},
	}
}
