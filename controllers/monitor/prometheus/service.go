package prometheus

import (
	"fmt"
	cachev1alpha1 "pulsar-operator/api/v1alpha1"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeService(c *cachev1alpha1.PulsarCluster) *v1.Service {
	var serviceType v1.ServiceType
	if c.Spec.Monitor.Prometheus.NodePort == 0 {
		serviceType = v1.ServiceTypeClusterIP
	} else {
		serviceType = v1.ServiceTypeNodePort
	}

	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeServiceName(c),
			Namespace: c.Namespace,
			Labels:    cachev1alpha1.MakeAllLabels(c, cachev1alpha1.MonitorComponent, cachev1alpha1.MonitorPrometheusComponent),
		},
		Spec: v1.ServiceSpec{
			Ports:    makeServicePorts(c),
			Type:     serviceType,
			Selector: cachev1alpha1.MakeAllLabels(c, cachev1alpha1.MonitorComponent, cachev1alpha1.MonitorPrometheusComponent),
		},
	}
}

func MakeServiceName(c *cachev1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-monitor-prometheus-service", c.GetName())
}

func makeServicePorts(c *cachev1alpha1.PulsarCluster) []v1.ServicePort {
	if c.Spec.Monitor.Prometheus.NodePort == 0 {
		return []v1.ServicePort{
			{
				Name: "prometheus",
				Port: cachev1alpha1.PulsarPrometheusServerPort,
			},
		}
	} else {
		return []v1.ServicePort{
			{
				Name:     "prometheus",
				NodePort: c.Spec.Monitor.Prometheus.NodePort,
				Port:     cachev1alpha1.PulsarPrometheusServerPort,
			},
		}
	}
}
