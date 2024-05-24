package toolset

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeDeployment(c *v1alpha1.PulsarCluster) *appsv1.Deployment {
	var replicas = c.Spec.Toolset.Replicas
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", c.Name, v1alpha1.ToolsetComponent),
			Namespace: c.Namespace,
			Labels:    v1alpha1.MakeComponentLabels(c, v1alpha1.ToolsetComponent),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: v1alpha1.MakeComponentLabels(c, v1alpha1.ToolsetComponent),
			},
			Template: makePodTemplate(c),
		},
	}
}
