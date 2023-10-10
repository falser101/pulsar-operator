package grafana

import (
	"fmt"
	"pulsar-operator/pkg/api/v1alpha1"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeService(c *v1alpha1.PulsarCluster) *v1.Service {
	var serviceType v1.ServiceType
	if c.Spec.Monitor.Grafana.NodePort == 0 {
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
			Labels:    v1alpha1.MakeAllLabels(c, v1alpha1.MonitorComponent, v1alpha1.MonitorGrafanaComponent),
		},
		Spec: v1.ServiceSpec{
			Ports:    makeServicePorts(c),
			Type:     serviceType,
			Selector: v1alpha1.MakeAllLabels(c, v1alpha1.MonitorComponent, v1alpha1.MonitorGrafanaComponent),
		},
	}
}

func MakeServiceName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-monitor-grafana-service", c.GetName())
}

func makeServicePorts(c *v1alpha1.PulsarCluster) []v1.ServicePort {
	if c.Spec.Monitor.Grafana.NodePort == 0 {
		return []v1.ServicePort{
			{
				Name: "grafana",
				Port: v1alpha1.PulsarGrafanaServerPort,
			},
		}
	} else {
		return []v1.ServicePort{
			{
				Name:     "grafana",
				NodePort: c.Spec.Monitor.Grafana.NodePort,
				Port:     v1alpha1.PulsarGrafanaServerPort,
			},
		}
	}
}
