package broker

import (
	"fmt"
	cachev1alpha1 "pulsar-operator/api/v1alpha1"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeService(c *cachev1alpha1.PulsarCluster) *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeServiceName(c),
			Namespace: c.Namespace,
			Labels:    cachev1alpha1.MakeComponentLabels(c, cachev1alpha1.BrokerComponent),
		},
		Spec: v1.ServiceSpec{
			Ports:     makeServicePorts(c),
			ClusterIP: v1.ClusterIPNone,
			Selector:  cachev1alpha1.MakeComponentLabels(c, cachev1alpha1.BrokerComponent),
		},
	}
}

func MakeServiceName(c *cachev1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-broker-service", c.GetName())
}

func makeServicePorts(c *cachev1alpha1.PulsarCluster) []v1.ServicePort {
	return []v1.ServicePort{
		{
			Name: "http",
			Port: cachev1alpha1.PulsarBrokerHttpServerPort,
		},
		{
			Name: "pulsar",
			Port: cachev1alpha1.PulsarBrokerPulsarServerPort,
		},
	}
}
