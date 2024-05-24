package proxy

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func MakeService(c *v1alpha1.PulsarCluster) *v1.Service {
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

func MakeServiceName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-%s", c.Name, v1alpha1.ProxyComponent)
}

func makeServicePorts(c *v1alpha1.PulsarCluster) []v1.ServicePort {
	var httpServerPort, pulsarServerPort int32
	if c.Spec.Proxy.HttpServerPort != 0 {
		httpServerPort = c.Spec.Proxy.HttpServerPort
	}
	if c.Spec.Proxy.PulsarServerPort != 0 {
		pulsarServerPort = c.Spec.Proxy.PulsarServerPort
	}
	return []v1.ServicePort{
		{
			Name: "http",
			Port: v1alpha1.PulsarBrokerHttpServerPort,
			TargetPort: intstr.IntOrString{
				IntVal: v1alpha1.PulsarBrokerHttpServerPort,
			},
			NodePort: httpServerPort,
		},
		{
			Name: "pulsar",
			Port: v1alpha1.PulsarBrokerPulsarServerPort,
			TargetPort: intstr.IntOrString{
				IntVal: v1alpha1.PulsarBrokerPulsarServerPort,
			},
			NodePort: pulsarServerPort,
		},
	}
}
