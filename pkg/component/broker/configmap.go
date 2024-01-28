package broker

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/pkg/component/zookeeper"
	v1 "k8s.io/api/core/v1"
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
		Data: makeConfigMapData(c),
	}
}

func makeConfigMapData(c *v1alpha1.Pulsar) (data map[string]string) {
	data = map[string]string{
		"PULSAR_MEM":                       PulsarMemData,
		"zookeeperServers":                 zookeeper.MakeServiceName(c),
		"configurationStoreServers":        zookeeper.MakeServiceName(c),
		"clusterName":                      c.GetName(),
		"managedLedgerDefaultEnsembleSize": ManagedLedgerDefaultEnsembleSize,
		"managedLedgerDefaultWriteQuorum":  ManagedLedgerDefaultWriteQuorum,
		"managedLedgerDefaultAckQuorum":    ManagedLedgerDefaultAckQuorum,
		"functionsWorkerEnabled":           FunctionsWorkerEnabled,
		"PF_pulsarFunctionsCluster":        c.GetName(),
	}
	if c.Spec.Authentication.Enabled {
		// broker.conf
		data["authenticationEnabled"] = "true"
		data["authenticationProviders"] = "org.apache.pulsar.broker.authentication.AuthenticationProviderToken"
		data["brokerClientAuthenticationPlugin"] = "org.apache.pulsar.client.impl.auth.AuthenticationToken"
		data["brokerClientAuthenticationParameters"] = "file:///pulsar/tokens/broker/token"
		data["tokenPublicKey"] = "file:///pulsar/keys/token/public.key"
		// client.conf
		data["authParams"] = "file:///pulsar/tokens/broker/token"
		data["authPlugin"] = "org.apache.pulsar.client.impl.auth.AuthenticationToken"
	}

	return
}

func MakeConfigMapName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-broker-configmap", c.GetName())
}
