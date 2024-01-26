package manager

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeServiceName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-manager-service", c.GetName())
}

func MakeService(c *v1alpha1.Pulsar) *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeServiceName(c),
			Namespace: c.Namespace,
			Labels:    v1alpha1.MakeComponentLabels(c, v1alpha1.ManagerComponent),
		},
		Spec: v1.ServiceSpec{
			Ports:    makeServicePorts(c),
			Type:     v1.ServiceTypeNodePort,
			Selector: v1alpha1.MakeComponentLabels(c, v1alpha1.ManagerComponent),
		},
	}
}

func makeServicePorts(c *v1alpha1.Pulsar) []v1.ServicePort {
	var frontendNodePort, backendNodePort int32
	if c.Spec.Manager.FrontendNodePort != 0 {
		frontendNodePort = c.Spec.Manager.FrontendNodePort
	}

	if c.Spec.Manager.BackendNodePort != 0 {
		backendNodePort = c.Spec.Manager.BackendNodePort
	}
	return []v1.ServicePort{
		{
			Name:     "frontend",
			NodePort: frontendNodePort,
			Port:     v1alpha1.PulsarManagerFrontendPort,
		},
		{
			Name:     "backend",
			NodePort: backendNodePort,
			Port:     v1alpha1.PulsarManagerBackendPort,
		},
	}
}
