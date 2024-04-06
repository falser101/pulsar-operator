package authentication

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeRole(c *v1alpha1.PulsarCluster) *v1.Role {
	return &v1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeRoleName(c),
			Namespace: c.Namespace,
		},
		Rules: []v1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"secrets", "namespaces"},
				Verbs:     []string{"get", "watch", "list", "create"},
			},
		},
	}
}

func MakeRoleName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-token-role", c.Name)
}
