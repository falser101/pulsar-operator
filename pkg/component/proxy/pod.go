package proxy

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/pkg/component/broker"
	"github.com/falser101/pulsar-operator/pkg/component/zookeeper"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func makePodSpec(c *v1alpha1.Pulsar) v1.PodSpec {
	return v1.PodSpec{
		Affinity: c.Spec.Broker.Pod.Affinity,
		InitContainers: []v1.Container{
			zookeeper.MakeWaitZookeeperReadyContainer(c),
			broker.MakeWaitBrokerReadyContainer(c),
		},
		Containers: []v1.Container{makeContainer(c)},
		Volumes:    makeVolumes(c),
	}
}

func makeContainer(c *v1alpha1.Pulsar) v1.Container {
	var volumeMounts = []v1.VolumeMount{
		{
			Name:      makeLog4j2Name(c),
			MountPath: "/pulsar/conf/log4j2.yaml",
			SubPath:   "log4j2.yaml",
		},
	}
	if c.Spec.Authentication.Enabled {
		volumeMounts = append(volumeMounts,
			v1.VolumeMount{
				Name:      "token-keys",
				MountPath: "/pulsar/keys",
				ReadOnly:  true,
			},
			v1.VolumeMount{
				Name:      "proxy-token",
				MountPath: "/pulsar/tokens",
				ReadOnly:  true,
			})
	}
	return v1.Container{
		Name:            "proxy",
		Image:           c.Spec.Proxy.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Proxy.Image.PullPolicy,
		Command:         makeContainerCommand(),
		Args:            makeContainerCommandArgs(),
		Ports:           makeContainerPort(c),
		EnvFrom:         makeContainerEnvFrom(c),
		ReadinessProbe: &v1.Probe{
			InitialDelaySeconds: 30,
			PeriodSeconds:       10,
			FailureThreshold:    10,
			SuccessThreshold:    1,
			TimeoutSeconds:      1,
			ProbeHandler: v1.ProbeHandler{
				HTTPGet: &v1.HTTPGetAction{
					Path: "/status.html",
					Port: intstr.IntOrString{
						IntVal: v1alpha1.PulsarBrokerHttpServerPort,
					},
					Scheme: v1.URISchemeHTTP,
				},
			},
		},
		LivenessProbe: &v1.Probe{
			InitialDelaySeconds: 30,
			PeriodSeconds:       10,
			FailureThreshold:    10,
			SuccessThreshold:    1,
			TimeoutSeconds:      1,
			ProbeHandler: v1.ProbeHandler{
				HTTPGet: &v1.HTTPGetAction{
					Path: "/status.html",
					Port: intstr.IntOrString{
						IntVal: v1alpha1.PulsarBrokerHttpServerPort,
					},
					Scheme: v1.URISchemeHTTP,
				},
			},
		},
		VolumeMounts: volumeMounts,
	}
}

func makeContainerCommand() []string {
	return []string{
		"sh",
		"-c",
	}
}

func makeContainerCommandArgs() []string {
	return []string{
		"bin/apply-config-from-env.py conf/proxy.conf; echo 'OK' > status; bin/pulsar proxy;",
	}
}

func makeContainerPort(c *v1alpha1.Pulsar) []v1.ContainerPort {
	return []v1.ContainerPort{
		{
			Name:          "http",
			ContainerPort: v1alpha1.PulsarBrokerHttpServerPort,
			Protocol:      v1.ProtocolTCP,
		},
		{
			Name:          "pulsar",
			ContainerPort: v1alpha1.PulsarBrokerPulsarServerPort,
			Protocol:      v1.ProtocolTCP,
		},
	}
}

func makeContainerEnvFrom(c *v1alpha1.Pulsar) []v1.EnvFromSource {
	froms := make([]v1.EnvFromSource, 0)

	var configRef v1.ConfigMapEnvSource
	configRef.Name = MakeConfigMapName(c)

	froms = append(froms, v1.EnvFromSource{ConfigMapRef: &configRef})
	return froms
}

func makeLog4j2Name(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-proxy-log4j2", c.Namespace)
}

func makeVolumes(c *v1alpha1.Pulsar) []v1.Volume {
	var defaultMode int32 = 420
	var volumes = []v1.Volume{
		{
			Name: makeLog4j2Name(c),
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{Name: MakeConfigMapName(c)},
					DefaultMode:          &defaultMode,
				},
			},
		},
	}
	if c.Spec.Authentication.Enabled {
		volumes = append(volumes,
			v1.Volume{
				Name: "token-keys",
				VolumeSource: v1.VolumeSource{
					Secret: &v1.SecretVolumeSource{
						SecretName: fmt.Sprintf("%s-token-asymmetric-key", c.Name),
						Items: []v1.KeyToPath{
							{
								Key:  "PUBLICKEY",
								Path: "token/public.key",
							},
						},
						DefaultMode: &defaultMode,
					},
				},
			},
			v1.Volume{
				Name: "proxy-token",
				VolumeSource: v1.VolumeSource{
					Secret: &v1.SecretVolumeSource{
						SecretName: fmt.Sprintf("%s-token-proxy-admin", c.Name),
						Items: []v1.KeyToPath{
							{
								Key:  "TOKEN",
								Path: "proxy/token",
							},
						},
						DefaultMode: &defaultMode,
					},
				},
			})
	}
	return volumes
}
