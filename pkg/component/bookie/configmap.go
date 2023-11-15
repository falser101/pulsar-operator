package bookie

import (
	"fmt"
	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/pkg/component/zookeeper"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeConfigMap(c *v1alpha1.Pulsar) *v1.ConfigMap {
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
			"BOOKIE_MEM":                MemData,
			"PULSAR_GC":                 PulsarGC,
			"PULSAR_MEM":                PulsarMem,
			"autoRecoveryDaemonEnabled": "false",
			"httpServerEnabled":         "true",
			"httpServerPort":            "8000",
			"journalDirectories":        "/pulsar/data/bookkeeper/journal",
			"journalMaxBackups":         "0",
			"ledgerDirectories":         "/pulsar/data/bookkeeper/ledgers",
			"zkServers":                 zookeeper.MakeServiceName(c) + ":2181",
			"zkLedgersRootPath":         "/ledgers",
			"statsProviderClass":        StatsProviderClass,
			"useHostNameAsBookieID":     "true",
		},
	}
}

func MakeConfigMapName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-bookie-configmap", c.GetName())
}
