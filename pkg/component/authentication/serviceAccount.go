package authentication

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeServiceAccount(c *v1alpha1.Pulsar) *v1.ServiceAccount {
	return &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeServiceAccountName(c),
			Namespace: c.Namespace,
		},
	}
}

func MakeServiceAccountName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-token-service-account", c.Name)
}
