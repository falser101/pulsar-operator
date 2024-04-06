package prometheus

import (
	"fmt"
	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func makePodSpec(c *v1alpha1.PulsarCluster) v1.PodSpec {
	pod := v1.PodSpec{
		Containers: []v1.Container{
			makeConfigMapReloadContainer(c),
			makeContainer(c),
		},
		Volumes:            makeVolumes(c),
		ServiceAccountName: MakeServiceAccountName(c),
	}
	if isUseEmptyDirVolume(c) {
		pod.Volumes = makeEmptyDirVolume(c)
	}
	return pod
}

func makeConfigMapReloadContainer(c *v1alpha1.PulsarCluster) v1.Container {
	return v1.Container{
		Name:                     makeConfigMapReloadContainerName(c),
		Image:                    "jimmidyson/configmap-reload:v0.3.0",
		ImagePullPolicy:          v1.PullIfNotPresent,
		Args:                     makeConfigMapReloadContainerArgs(c),
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: v1.TerminationMessageReadFile,
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      "config-volume",
				ReadOnly:  true,
				MountPath: "/etc/config",
			},
		},
	}
}

func makeConfigMapReloadContainerArgs(c *v1alpha1.PulsarCluster) []string {
	return []string{
		"--volume-dir=/etc/config",
		"--webhook-url=http://127.0.0.1:9090/-/reload",
	}
}

func makeConfigMapReloadContainerName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-prometheus-configmap-reload", c.Name)
}

func makeContainer(c *v1alpha1.PulsarCluster) v1.Container {
	return v1.Container{
		Name:            "prometheus",
		Image:           c.Spec.Monitor.Prometheus.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Monitor.Prometheus.Image.PullPolicy,
		Args: []string{
			"--config.file=/etc/config/prometheus.yml",
			"--storage.tsdb.retention.time=15d",
			"--storage.tsdb.path=/prometheus",
			"--web.console.libraries=/etc/prometheus/console_libraries",
			"--web.console.templates=/etc/prometheus/consoles",
			"--web.enable-lifecycle",
		},
		Ports: makeContainerPort(c),
		LivenessProbe: &v1.Probe{
			FailureThreshold:    10,
			InitialDelaySeconds: 30,
			PeriodSeconds:       10,
			SuccessThreshold:    1,
			TimeoutSeconds:      1,
			ProbeHandler: v1.ProbeHandler{
				HTTPGet: &v1.HTTPGetAction{
					Path: "/-/healthy",
					Port: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 9090,
					},
					Scheme: v1.URISchemeHTTP,
				},
			},
		},
		ReadinessProbe: &v1.Probe{
			FailureThreshold:    10,
			InitialDelaySeconds: 30,
			TimeoutSeconds:      1,
			PeriodSeconds:       10,
			SuccessThreshold:    1,
			ProbeHandler: v1.ProbeHandler{
				HTTPGet: &v1.HTTPGetAction{
					Path: "/-/ready",
					Port: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 9090,
					},
					Scheme: v1.URISchemeHTTP,
				},
			},
		},
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: v1.TerminationMessageReadFile,
		VolumeMounts:             makeContainerVolumeMount(c),
	}
}

func makeContainerPort(c *v1alpha1.PulsarCluster) []v1.ContainerPort {
	return []v1.ContainerPort{
		{
			Name:          "prometheus",
			ContainerPort: v1alpha1.PulsarPrometheusServerPort,
			Protocol:      v1.ProtocolTCP,
		},
	}
}

func makeContainerVolumeMount(c *v1alpha1.PulsarCluster) []v1.VolumeMount {
	return []v1.VolumeMount{
		{
			Name:      ConfigVolumeName,
			MountPath: "/etc/config",
		},
		{
			Name:      DataVolumeName,
			MountPath: "/prometheus",
		},
		{
			Name:      ClientTokenName,
			MountPath: "/pulsar/tokens",
			ReadOnly:  true,
		},
	}
}

func makeVolumes(c *v1alpha1.PulsarCluster) []v1.Volume {
	return []v1.Volume{
		{
			Name:         DataVolumeName,
			VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}},
		},
		{
			Name:         ConfigVolumeName,
			VolumeSource: v1.VolumeSource{ConfigMap: &v1.ConfigMapVolumeSource{LocalObjectReference: v1.LocalObjectReference{Name: MakeConfigMapName(c)}}},
		},
		{
			Name: ClientTokenName,
			VolumeSource: v1.VolumeSource{Secret: &v1.SecretVolumeSource{
				SecretName: fmt.Sprintf("%s-token-admin", c.Name),
			}},
		},
	}
}

func isUseEmptyDirVolume(c *v1alpha1.PulsarCluster) bool {
	return c.Spec.Monitor.Prometheus.StorageClassName == ""
}
