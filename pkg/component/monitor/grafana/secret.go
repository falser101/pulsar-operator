package grafana

import (
	"fmt"
	"github.com/falser101/pulsar-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeSecret(c *v1alpha1.Pulsar) *v1.Secret {
	return &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      makeGrafanaSecretName(c),
			Namespace: c.Namespace,
		},
		Data: map[string][]byte{
			"GRAFANA_ADMIN_PASSWORD": []byte(c.Spec.Monitor.Grafana.Security.AdminPassword),
			"GRAFANA_ADMIN_USER":     []byte(c.Spec.Monitor.Grafana.Security.AdminUser),
		},
		Type: v1.SecretTypeOpaque,
	}
}

func makeGrafanaSecretName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-monitor-grafana-secret", c.Name)
}
