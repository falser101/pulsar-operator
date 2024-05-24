package toolset

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/internal/component/proxy"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeConfigMap(c *v1alpha1.PulsarCluster) *v1.ConfigMap {
	var configData = map[string]string{
		"BOOKIE_LOG_APPENDER": "RollingFile",
		"webServiceUrl":       fmt.Sprintf("http://%s:8080", proxy.MakeServiceName(c)),
		"brokerServiceUrl":    fmt.Sprintf("pulsar://%s:6650", proxy.MakeServiceName(c)),
	}
	if c.Spec.Auth.AuthenticationEnabled {
		configData["authParams"] = "file:///pulsar/tokens/client/token"
		configData["authPlugin"] = "org.apache.pulsar.client.impl.auth.AuthenticationToken"
	}
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeConfigMapName(c),
			Namespace: c.Namespace,
		},
		Data: configData,
	}
}

func MakeConfigMapName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-%s", c.Name, v1alpha1.ToolsetComponent)
}
