package authentication

import (
	"github.com/falser101/pulsar-operator/api/v1alpha1"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeRoleBinding(c *v1alpha1.Pulsar) *v1.RoleBinding {
	return &v1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeRoleBindingName(c),
			Namespace: c.Namespace,
		},
		Subjects: []v1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      c.Name,
				Namespace: c.Namespace,
			},
		},
		RoleRef: v1.RoleRef{
			Kind:     "Role",
			Name:     MakeRoleName(c),
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
}

func MakeRoleBindingName(c *v1alpha1.Pulsar) string {
	return c.Name + "-role-binding"
}
