package zookeeper

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeServiceAccount(c *v1alpha1.PulsarCluster) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      makeSerserviceAccountName(c),
			Namespace: c.Namespace,
			Labels:    v1alpha1.MakeComponentLabels(c, v1alpha1.ZookeeperComponent),
		},
	}
}

func makeSerserviceAccountName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-%s", c.Name, v1alpha1.ZookeeperComponent)
}
