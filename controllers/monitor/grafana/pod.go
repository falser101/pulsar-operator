package grafana

import (
	"fmt"
	cachev1alpha1 "pulsar-operator/api/v1alpha1"
	"pulsar-operator/controllers/monitor/prometheus"

	"k8s.io/api/core/v1"
)

func makePodSpec(c *cachev1alpha1.PulsarCluster) v1.PodSpec {
	return v1.PodSpec{
		Containers: []v1.Container{makeContainer(c)},
	}
}

func makeContainer(c *cachev1alpha1.PulsarCluster) v1.Container {
	return v1.Container{
		Name:            "grafana",
		Image:           cachev1alpha1.MonitorGrafanaImage,
		ImagePullPolicy: cachev1alpha1.DefaultContainerPolicy,
		Ports:           makeContainerPort(c),
		Env:             makeContainerEnv(c),
	}
}

func makeContainerPort(c *cachev1alpha1.PulsarCluster) []v1.ContainerPort {
	return []v1.ContainerPort{
		{
			Name:          "grafana",
			ContainerPort: cachev1alpha1.PulsarGrafanaServerPort,
			Protocol:      v1.ProtocolTCP,
		},
	}
}

func makeContainerEnv(c *cachev1alpha1.PulsarCluster) []v1.EnvVar {
	prometheusUrl := fmt.Sprintf("http://%s:%d/", prometheus.MakeServiceName(c), cachev1alpha1.PulsarPrometheusServerPort)
	env := make([]v1.EnvVar, 0)
	p := v1.EnvVar{
		Name:  "PROMETHEUS_URL",
		Value: prometheusUrl,
	}
	env = append(env, p)
	return env
}
