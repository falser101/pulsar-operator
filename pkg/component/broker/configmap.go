package broker

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
			"PULSAR_MEM":                       PulsarMemData,
			"zookeeperServers":                 zookeeper.MakeServiceName(c),
			"configurationStoreServers":        zookeeper.MakeServiceName(c),
			"clusterName":                      c.GetName(),
			"managedLedgerDefaultEnsembleSize": ManagedLedgerDefaultEnsembleSize,
			"managedLedgerDefaultWriteQuorum":  ManagedLedgerDefaultWriteQuorum,
			"managedLedgerDefaultAckQuorum":    ManagedLedgerDefaultAckQuorum,
			"functionsWorkerEnabled":           FunctionsWorkerEnabled,
			"PF_pulsarFunctionsCluster":        c.GetName(),
		},
	}
}

func MakeConfigMapName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-broker-configmap", c.GetName())
}
