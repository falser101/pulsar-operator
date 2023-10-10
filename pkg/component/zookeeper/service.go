package zookeeper

import (
	"fmt"
	"pulsar-operator/pkg/api/v1alpha1"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeService(c *v1alpha1.PulsarCluster) *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        MakeServiceName(c),
			Namespace:   c.Namespace,
			Labels:      v1alpha1.MakeComponentLabels(c, v1alpha1.ZookeeperComponent),
			Annotations: ServiceAnnotations,
		},
		Spec: v1.ServiceSpec{
			Ports:     makeServicePorts(c),
			ClusterIP: v1.ClusterIPNone,
			Selector:  v1alpha1.MakeComponentLabels(c, v1alpha1.ZookeeperComponent),
		},
	}
}

func MakeServiceName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-zookeeper-service", c.GetName())
}

func makeServicePorts(c *v1alpha1.PulsarCluster) []v1.ServicePort {
	return []v1.ServicePort{
		{
			Name: "server",
			Port: v1alpha1.ZookeeperContainerServerDefaultPort,
		},
		{
			Name: "leader-election",
			Port: v1alpha1.ZookeeperContainerLeaderElectionPort,
		},
	}
}
