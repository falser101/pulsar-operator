package prometheus

import (
	"fmt"
	"github.com/falser101/pulsar-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeClusterRole(c *v1alpha1.Pulsar) *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   MakeClusterRoleName(c),
			Labels: v1alpha1.MakeAllLabels(c, v1alpha1.MonitorComponent, v1alpha1.MonitorPrometheusComponent),
		},
		Rules: makeClusterRoleRules(c),
	}
}

func MakeClusterRoleName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-prometheus-cluster-role", c.GetName())
}

func makeClusterRoleRules(c *v1alpha1.Pulsar) []rbacv1.PolicyRule {
	result := make([]rbacv1.PolicyRule, 0)

	rule1 := rbacv1.PolicyRule{
		APIGroups: []string{""},
		Resources: []string{
			"nodes",
			"nodes/proxy",
			"services",
			"endpoints",
			"pods",
		},
		Verbs: []string{
			"get",
			"list",
			"watch",
		},
	}
	result = append(result, rule1)

	rule2 := rbacv1.PolicyRule{
		NonResourceURLs: []string{
			"/metrics",
		},
		Verbs: []string{
			"get",
		},
	}
	result = append(result, rule2)
	return result
}

func MakeServiceAccount(c *v1alpha1.Pulsar) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeServiceAccountName(c),
			Namespace: c.Namespace,
			Labels:    v1alpha1.MakeAllLabels(c, v1alpha1.MonitorComponent, v1alpha1.MonitorPrometheusComponent),
		},
	}
}

func MakeServiceAccountName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-prometheus-serviceaccount", c.GetName())
}

func MakeClusterRoleBinding(c *v1alpha1.Pulsar) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   MakeClusterRoleBindingName(c),
			Labels: v1alpha1.MakeAllLabels(c, v1alpha1.MonitorComponent, v1alpha1.MonitorPrometheusComponent),
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     MakeClusterRoleName(c),
		},
		Subjects: makeSubjects(c),
	}
}

func MakeClusterRoleBindingName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-prometheus-clusterrole-binding", c.GetName())
}

func makeSubjects(c *v1alpha1.Pulsar) []rbacv1.Subject {
	result := make([]rbacv1.Subject, 0)

	s := rbacv1.Subject{
		Kind:      "ServiceAccount",
		Name:      MakeServiceAccountName(c),
		Namespace: c.GetNamespace(),
	}
	result = append(result, s)
	return result
}
