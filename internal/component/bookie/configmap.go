package bookie

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeConfigMap(c *v1alpha1.PulsarCluster) *v1.ConfigMap {
	var data = map[string]string{
		"autoRecoveryDaemonEnabled":        "false",
		"journalMaxBackups":                "0",
		"journalDirectories":               "/pulsar/data/bookkeeper/journal",
		"PULSAR_PREFIX_journalDirectories": "/pulsar/data/bookkeeper/journal",
		"ledgerDirectories":                "/pulsar/data/bookkeeper/ledgers",
	}
	for key, value := range c.Spec.Bookie.ConfigData {
		data[key] = value
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

func MakeConfigMapName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-%s", c.Name, v1alpha1.BookieComponent)
}
