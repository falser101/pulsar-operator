package broker

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/internal/component/zookeeper"
	v1 "k8s.io/api/core/v1"
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
		Data: makeConfigMapData(c),
	}
}

func makeConfigMapData(c *v1alpha1.PulsarCluster) (data map[string]string) {
	data = map[string]string{
		// metadata setting
		"zookeeperServers":          zookeeper.Connect(c),
		"configurationStoreServers": zookeeper.Connect(c),
		// Broker settings
		"clusterName":                         c.Name,
		"exposeTopicLevelMetricsInPrometheus": "true",
		"numHttpServerThreads":                "8",
		"zooKeeperSessionTimeoutMillis":       "30000",
		"statusFilePath":                      "/pulsar/logs/status",
		"functionsWorkerEnabled":              FunctionsWorkerEnabled,

		"webServicePort":    c.Spec.Broker.Ports.Http,
		"brokerServicePort": c.Spec.Broker.Ports.Pulsar,
	}
	if c.Spec.Auth.AuthorizationEnabled {
		data["authorizationEnabled"] = "true"
	}
	if c.Spec.Auth.AuthenticationEnabled {
		// broker.conf
		data["authenticationEnabled"] = "true"
		data["authenticationProviders"] = "org.apache.pulsar.broker.authentication.AuthenticationProviderToken"
		data["brokerClientAuthenticationPlugin"] = "org.apache.pulsar.impl.auth.AuthenticationToken"
		data["brokerClientAuthenticationParameters"] = "file:///pulsar/tokens/broker/token"
		data["tokenPublicKey"] = "file:///pulsar/keys/token/public.key"
		// client.conf
		data["authParams"] = "file:///pulsar/tokens/broker/token"
		data["authPlugin"] = "org.apache.pulsar.impl.auth.AuthenticationToken"
	}
	for key, value := range c.Spec.Broker.ConfigData {
		data[key] = value
	}
	return
}

func MakeConfigMapName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-%s", c.Name, v1alpha1.BrokerComponent)
}
