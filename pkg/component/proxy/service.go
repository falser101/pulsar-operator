package proxy

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func MakeService(c *v1alpha1.Pulsar) *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeServiceName(c),
			Namespace: c.Namespace,
			Labels:    v1alpha1.MakeComponentLabels(c, v1alpha1.ProxyComponent),
		},
		Spec: v1.ServiceSpec{
			Ports:    makeServicePorts(c),
			Type:     v1.ServiceTypeLoadBalancer,
			Selector: v1alpha1.MakeComponentLabels(c, v1alpha1.ProxyComponent),
		},
	}
}

func MakeServiceName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-proxy-service", c.GetName())
}

func makeServicePorts(c *v1alpha1.Pulsar) []v1.ServicePort {
	return []v1.ServicePort{
		{
			Name: "http",
			Port: v1alpha1.PulsarBrokerHttpServerPort,
			TargetPort: intstr.IntOrString{
				IntVal: v1alpha1.PulsarBrokerHttpServerPort,
			},
			NodePort: 30080,
		},
		{
			Name: "pulsar",
			Port: v1alpha1.PulsarBrokerPulsarServerPort,
			TargetPort: intstr.IntOrString{
				IntVal: v1alpha1.PulsarBrokerPulsarServerPort,
			},
			NodePort: 30650,
		},
	}
}
