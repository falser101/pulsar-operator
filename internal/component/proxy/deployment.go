package proxy

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeDeployment(c *v1alpha1.PulsarCluster) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      makeDeploymentName(c),
			Namespace: c.Namespace,
			Labels:    v1alpha1.MakeComponentLabels(c, v1alpha1.ProxyComponent),
		},
		Spec: makeDeploymentSpec(c),
	}
}

func makeDeploymentName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-%s", c.Name, v1alpha1.ProxyComponent)
}

func makeDeploymentSpec(c *v1alpha1.PulsarCluster) appsv1.DeploymentSpec {
	return appsv1.DeploymentSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: v1alpha1.MakeComponentLabels(c, v1alpha1.ProxyComponent),
		},
		Replicas: &c.Spec.Proxy.Replicas,
		Template: makeDeploymentPodTemplate(c),
	}
}

func makeDeploymentPodTemplate(c *v1alpha1.PulsarCluster) v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: c.Name,
			Labels:       v1alpha1.MakeComponentLabels(c, v1alpha1.ProxyComponent),
			Annotations:  DeploymentAnnotations,
		},
		Spec: makePodSpec(c),
	}
}