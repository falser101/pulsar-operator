package zookeeper

import (
	"fmt"
	"github.com/falser101/pulsar-operator/api/v1alpha1"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeConfigMap(c *v1alpha1.PulsarCluster) *v1.ConfigMap {
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeConfigMapName(c),
			Namespace: c.Namespace,
		},
		Data: map[string]string{
			"PULSAR_MEM": PulsarMemData,
			"PULSAR_GC":  PulsarGCData,
		},
	}
}

func MakeConfigMapName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-zookeeper-configmap", c.GetName())
}