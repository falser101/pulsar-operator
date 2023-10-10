package manager

import (
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cachev1alpha1 "pulsar-operator/pkg/api/v1alpha1"
	"pulsar-operator/pkg/component/bookie"
	"pulsar-operator/pkg/component/broker"
)

func MakeConfigMap(c *cachev1alpha1.PulsarCluster) *v1.ConfigMap {
	data := make(map[string]string)
	if c.Spec.Manager.ConfigMap == nil {
		data[BackendEntrypointKey] = MakeBackendEntrypoint(c)
		data[EntrypointKey] = EntrypointValue
	}
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeConfigMapName(c),
			Namespace: c.Namespace,
		},
		Data: data,
	}
}

func MakeConfigMapName(c *cachev1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-manager-configmap", c.GetName())
}

func MakeBackendEntrypoint(c *cachev1alpha1.PulsarCluster) string {
	return fmt.Sprintf(BackendEntrypointValue, bookie.MakeServiceName(c), c.Name, broker.MakeServiceName(c))
}
