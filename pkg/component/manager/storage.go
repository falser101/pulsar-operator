package manager

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makeManagerTokensName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-manager-tokens", c.Name)
}

func makeManagerTokenKeysName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-manager-token-keys", c.Name)
}

func makeManagerBackendScriptName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-manager-backend-script", c.Name)
}

func makeManagerScriptName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-manager-script", c.Name)
}

func makeManagerDataName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-manager-data", c.Name)
}

func makeEmptyDirVolume(c *v1alpha1.Pulsar) []v1.Volume {
	var scriptMode int32 = 493
	var tokenMode int32 = 420
	var volumes = []v1.Volume{
		{
			Name:         makeManagerDataName(c),
			VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}},
		},
		{
			Name: makeManagerScriptName(c),
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{Name: MakeConfigMapName(c)},
					DefaultMode:          &scriptMode,
				}},
		},
		{
			Name: makeManagerBackendScriptName(c),
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{Name: MakeConfigMapName(c)},
					DefaultMode:          &scriptMode,
				}},
		},
	}
	if c.Spec.Authentication.Enabled {
		volumes = append(volumes, v1.Volume{
			Name: makeManagerTokenKeysName(c),
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: makeAsymmetricKeyName(c),
					Items: []v1.KeyToPath{
						{
							Key:  "PUBLICKEY",
							Path: "token/public.key",
						},
						{
							Key:  "PRIVATEKEY",
							Path: "token/private.key",
						},
					},
					DefaultMode: &tokenMode,
				}}},
			v1.Volume{
				Name: makeManagerTokensName(c),
				VolumeSource: v1.VolumeSource{
					Secret: &v1.SecretVolumeSource{
						SecretName: makeManagerAdminSecretName(c),
						Items: []v1.KeyToPath{
							{
								Key:  "TOKEN",
								Path: "pulsar_manager/token",
							},
						},
						DefaultMode: &tokenMode,
					}},
			})
	}

	return volumes
}

func makeManagerAdminSecretName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-token-console-admin", c.Name)
}

func makeAsymmetricKeyName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-token-asymmetric-key", c.Name)
}

// PV/PVC
func makeVolumeClaimTemplates(c *v1alpha1.Pulsar) []v1.PersistentVolumeClaim {
	return []v1.PersistentVolumeClaim{
		makeManagerDataVolumeClaimTemplate(c),
	}
}

func makeManagerDataVolumeClaimTemplate(c *v1alpha1.Pulsar) v1.PersistentVolumeClaim {
	return v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      makeManagerDataName(c),
			Namespace: c.Namespace,
		},
		Spec: makeManagerDataVolumeClaimSpec(c),
	}
}

func makeManagerDataVolumeClaimSpec(c *v1alpha1.Pulsar) v1.PersistentVolumeClaimSpec {
	capacity := fmt.Sprintf("%dGi", c.Spec.Manager.StorageCapacity)
	return v1.PersistentVolumeClaimSpec{
		AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
		Resources:        v1.ResourceRequirements{Requests: v1.ResourceList{v1.ResourceStorage: resource.MustParse(capacity)}},
		StorageClassName: &c.Spec.Manager.StorageClassName,
	}
}
