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
	var servicePorts = make([]v1.ServicePort, 0, 2)
	if c.Spec.Manager.FrontendNodePort == 0 {
		servicePorts = append(servicePorts, v1.ServicePort{
			Name:     "frontend",
			NodePort: v1alpha1.PulsarManagerFrontendNodePort,
			Port:     v1alpha1.PulsarManagerFrontendPort,
		})
	} else {
		servicePorts = append(servicePorts, v1.ServicePort{
			Name:     "frontend",
			NodePort: c.Spec.Manager.FrontendNodePort,
			Port:     v1alpha1.PulsarManagerFrontendPort,
		})
	}

	if c.Spec.Manager.BackendNodePort == 0 {
		servicePorts = append(servicePorts, v1.ServicePort{
			Name:     "backend",
			NodePort: v1alpha1.PulsarManagerBackNodePort,
			Port:     v1alpha1.PulsarManagerBackendPort,
		})
	} else {
		servicePorts = append(servicePorts, v1.ServicePort{
			Name:     "backend",
			NodePort: c.Spec.Manager.BackendNodePort,
			Port:     v1alpha1.PulsarManagerBackendPort,
		})
	}
	return servicePorts
}
