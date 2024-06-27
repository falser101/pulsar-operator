package zookeeper

import (
	"github.com/falser101/pulsar-operator/api/v1alpha1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeConfigMap(c *v1alpha1.PulsarCluster) *v1.ConfigMap {
	var data = map[string]string{
		"dataDir":                         "/pulsar/data/zookeeper",
		"PULSAR_PREFIX_serverCnxnFactory": "org.apache.zookeeper.server.NIOServerCnxnFactory",
		"serverCnxnFactory":               "org.apache.zookeeper.server.NIOServerCnxnFactory",
	}
	for key, val := range c.Spec.Zookeeper.ConfigData {
		data[key] = val
	}
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeName(c),
			Namespace: c.Namespace,
		},
		Data: data,
	}
}
