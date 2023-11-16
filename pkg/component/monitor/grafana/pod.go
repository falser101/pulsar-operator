package grafana

import (
	"fmt"
	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/pkg/component/monitor/prometheus"
	"k8s.io/api/core/v1"
)

func makePodSpec(c *v1alpha1.Pulsar) v1.PodSpec {
	return v1.PodSpec{
		Containers: []v1.Container{makeContainer(c)},
		Volumes:    makeVolumes(c),
	}
}

func makeVolumes(c *v1alpha1.Pulsar) []v1.Volume {
	var defaultMode int32 = 420
	var dataVolume v1.Volume
	if isUseEmptyDirVolume(c) {
		dataVolume = makeEmptyDirDataVolume(c)
	} else {
		dataVolume = makePVCDataVolume(c)
	}
	return []v1.Volume{
		{
			Name: "cfg",
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: MakeConfigMapName(c),
					},
					DefaultMode: &defaultMode,
				},
			},
		},
		dataVolume,
	}
}

func makePVCDataVolume(c *v1alpha1.Pulsar) v1.Volume {
	return v1.Volume{
		Name: makeDataVolumeName(c),
		VolumeSource: v1.VolumeSource{
			PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
				ClaimName: makePVCName(c),
			},
		},
	}
}

func makeContainer(c *v1alpha1.Pulsar) v1.Container {
	return v1.Container{
		Name:            "grafana",
		Image:           c.Spec.Monitor.Grafana.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Monitor.Grafana.Image.PullPolicy,
		Ports:           makeContainerPort(c),
		Env:             makeContainerEnv(c),
		VolumeMounts:    makeVolumeMounts(c),
	}
}

func makeVolumeMounts(c *v1alpha1.Pulsar) []v1.VolumeMount {
	return []v1.VolumeMount{
		{
			Name:      "cfg",
			MountPath: "/pulsar/conf/grafana.ini",
			SubPath:   "grafana.ini",
		},
		{
			Name:      makeDataVolumeName(c),
			MountPath: "/var/lib/grafana/pulsar",
		},
	}
}

func makeDataVolumeName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-monitor-grafana-data", c.Name)
}

func makeContainerPort(c *v1alpha1.Pulsar) []v1.ContainerPort {
	return []v1.ContainerPort{
		{
			Name:          "grafana",
			ContainerPort: v1alpha1.PulsarGrafanaServerPort,
			Protocol:      v1.ProtocolTCP,
		},
	}
}

func makeContainerEnv(c *v1alpha1.Pulsar) []v1.EnvVar {
	prometheusUrl := fmt.Sprintf("http://%s:%d/", prometheus.MakeServiceName(c), v1alpha1.PulsarPrometheusServerPort)
	env := []v1.EnvVar{
		{
			Name:  "PROMETHEUS_URL",
			Value: prometheusUrl,
		},
		{
			Name:  "PULSAR_PROMETHEUS_URL",
			Value: prometheusUrl,
		},
		{
			Name:  "PULSAR_CLUSTER",
			Value: c.Name,
		},
		{
			Name: "GRAFANA_ADMIN_USER",
			ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: makeGrafanaSecretName(c),
				},
				Key: "GRAFANA_ADMIN_USER",
			}},
		},
		{
			Name: "GRAFANA_ADMIN_PASSWORD",
			ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: makeGrafanaSecretName(c),
				},
				Key: "GRAFANA_ADMIN_PASSWORD",
			}},
		},
		{
			Name:  "GRAFANA_CFG_FILE",
			Value: "/pulsar/conf/grafana.ini",
		},
		{
			Name:  "GF_PATHS_DATA",
			Value: "/var/lib/grafana/pulsar/data",
		},
		{
			Name:  "GF_PATHS_PLUGINS",
			Value: "/var/lib/grafana/pulsar/plugin",
		},
		{
			Name:  "GF_PATHS_PROVISIONING",
			Value: "/var/lib/grafana/pulsar_provisioning",
		},
		{
			Name:  "GRAFANA_ROOT_URL",
			Value: "/grafana/",
		},
		{
			Name:  "GRAFANA_SERVE_FROM_SUB_PATH",
			Value: "true",
		},
	}
	return env
}

// EmptyDir volume
func makeEmptyDirDataVolume(c *v1alpha1.Pulsar) v1.Volume {
	return v1.Volume{
		Name:         makeDataVolumeName(c),
		VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}},
	}
}

func isUseEmptyDirVolume(c *v1alpha1.Pulsar) bool {
	return c.Spec.Monitor.Grafana.StorageClassName == ""
}
