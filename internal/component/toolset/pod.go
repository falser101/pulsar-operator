package toolset

import (
	"github.com/falser101/pulsar-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makePodTemplate(c *v1alpha1.PulsarCluster) v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: v1alpha1.MakeComponentLabels(c, v1alpha1.ToolsetComponent),
		},
		Spec: v1.PodSpec{
			Containers: makeContainers(c),
			Volumes:    makeVolumes(c),
		},
	}
}

func makeContainers(c *v1alpha1.PulsarCluster) []v1.Container {
	return []v1.Container{
		makePulsarContainer(c),
	}
}

func makePulsarContainer(c *v1alpha1.PulsarCluster) v1.Container {
	return v1.Container{
		Name:            "pulsar",
		Image:           c.Spec.Toolset.Image.GenerateImage(),
		ImagePullPolicy: v1.PullIfNotPresent,
		Command: []string{
			"sh",
			"-c",
		},
		Args:         makePulsarContainerArgs(c),
		EnvFrom:      makePulsarContainerEnvFrom(c),
		Resources:    makePulsarContainerResources(c),
		VolumeMounts: makePulsarContainerVolumeMounts(c),
	}
}

func makePulsarContainerArgs(c *v1alpha1.PulsarCluster) []string {
	return []string{
		`bin/apply-config-from-env.py conf/client.conf;
		bin/apply-config-from-env.py conf/bookkeeper.conf;
		echo "Configuring pulsarctl context ...";
		mkdir -p /root/.config/pulsar;
		cp /pulsar/conf/pulsarctl.config /root/.config/pulsar/config;
		echo "Successfully configured pulsarctl context."
		"sleep 10000000000`,
	}
}
func makePulsarContainerResources(c *v1alpha1.PulsarCluster) v1.ResourceRequirements {
	return v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("1"),
			v1.ResourceMemory: resource.MustParse("1Gi"),
		},
		Requests: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("100m"),
			v1.ResourceMemory: resource.MustParse("256Mi"),
		},
	}
}

func makePulsarContainerVolumeMounts(c *v1alpha1.PulsarCluster) []v1.VolumeMount {
	var VolumeMounts = []v1.VolumeMount{
		{
			Name:      "log4j2",
			MountPath: "/pulsar/conf/log4j2.yaml",
			SubPath:   "log4j2.yaml",
		},
		{
			Name:      "pulsarctl",
			MountPath: "/pulsar/conf/pulsarctl.config",
			SubPath:   "pulsarctl.config",
		},
	}
	if c.Spec.Auth.AuthenticationEnabled {
		VolumeMounts = append(VolumeMounts,
			v1.VolumeMount{
				Name:      "token-private-key",
				MountPath: "/pulsar/token-private-key",
				ReadOnly:  true,
			},
			v1.VolumeMount{
				Name:      "client-token",
				MountPath: "/pulsar/tokens",
				ReadOnly:  true,
			})
	}
	return VolumeMounts
}

func makePulsarContainerEnvFrom(c *v1alpha1.PulsarCluster) []v1.EnvFromSource {
	return []v1.EnvFromSource{
		{
			ConfigMapRef: &v1.ConfigMapEnvSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: MakeConfigMapName(c),
				},
			},
		},
	}
}

func makeVolumes(c *v1alpha1.PulsarCluster) []v1.Volume {
	return []v1.Volume{
		{
			Name: "token-private-key",
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: "test-token-asymmetric-key",
				},
			},
		},
		{
			Name: "client-token",
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: "test-client-token",
				},
			},
		},
		{
			Name: "log4j2",
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: MakeConfigMapName(c),
					},
				},
			},
		},
		{
			Name: "pulsarctl",
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: MakeConfigMapName(c),
					},
				},
			},
		},
	}
}
