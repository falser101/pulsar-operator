package bookie

import (
	"fmt"
	"github.com/falser101/pulsar-operator/api/v1alpha1"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makeJournalDataVolumeName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("journal-disk-volume-pvc")
}

func makeLedgersDataVolumeName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("ledgers-disk-volume-pvc")
}

// PV/PVC
func makeVolumeClaimTemplates(c *v1alpha1.Pulsar) []v1.PersistentVolumeClaim {
	return []v1.PersistentVolumeClaim{
		makeJournalDataVolumeClaimTemplate(c),
		makeLedgersDataVolumeClaimTemplate(c),
	}
}

func makeJournalDataVolumeClaimTemplate(c *v1alpha1.Pulsar) v1.PersistentVolumeClaim {
	return v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      makeJournalDataVolumeName(c),
			Namespace: c.Namespace,
		},
		Spec: makeJournalDataVolumeClaimSpec(c),
	}
}

func makeJournalDataVolumeClaimSpec(c *v1alpha1.Pulsar) v1.PersistentVolumeClaimSpec {
	capacity := fmt.Sprintf("%dGi", c.Spec.Bookie.JournalStorageCapacity)
	return v1.PersistentVolumeClaimSpec{
		AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
		Resources:        v1.ResourceRequirements{Requests: v1.ResourceList{v1.ResourceStorage: resource.MustParse(capacity)}},
		StorageClassName: &c.Spec.Bookie.StorageClassName,
	}
}

func makeLedgersDataVolumeClaimTemplate(c *v1alpha1.Pulsar) v1.PersistentVolumeClaim {
	return v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      makeLedgersDataVolumeName(c),
			Namespace: c.Namespace,
		},
		Spec: makeLedgersDataVolumeClaimSpec(c),
	}
}

func makeLedgersDataVolumeClaimSpec(c *v1alpha1.Pulsar) v1.PersistentVolumeClaimSpec {
	capacity := fmt.Sprintf("%dGi", c.Spec.Bookie.LedgersStorageCapacity)
	return v1.PersistentVolumeClaimSpec{
		AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
		Resources:        v1.ResourceRequirements{Requests: v1.ResourceList{v1.ResourceStorage: resource.MustParse(capacity)}},
		StorageClassName: &c.Spec.Bookie.StorageClassName,
	}
}

// EmptyDir volume
func makeEmptyDirVolume(c *v1alpha1.Pulsar) []v1.Volume {
	return []v1.Volume{
		{
			Name:         makeJournalDataVolumeName(c),
			VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}},
		},
		{
			Name:         makeLedgersDataVolumeName(c),
			VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}},
		},
	}
}
