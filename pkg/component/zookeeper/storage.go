package zookeeper

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makeZookeeperDateVolumeName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-zookeeper-data", c.Name)
}

func makeVolumeClaimTemplates(c *v1alpha1.Pulsar) []v1.PersistentVolumeClaim {
	return []v1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: makeZookeeperDateVolumeName(c),
			},
			Spec: makeZookeeperDataVolumeClaimSpec(c),
		},
	}
}

func makeZookeeperDataVolumeClaimSpec(c *v1alpha1.Pulsar) v1.PersistentVolumeClaimSpec {
	capacity := fmt.Sprintf("%dGi", c.Spec.Zookeeper.StorageCapacity)
	return v1.PersistentVolumeClaimSpec{
		AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
		Resources:        v1.ResourceRequirements{Requests: v1.ResourceList{v1.ResourceStorage: resource.MustParse(capacity)}},
		StorageClassName: &c.Spec.Zookeeper.StorageClassName,
	}
}

func makeEmptyDirVolume(c *v1alpha1.Pulsar) []v1.Volume {
	return []v1.Volume{
		{
			Name: makeZookeeperDateVolumeName(c),
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{},
			},
		},
	}
}
