package broker

import (
	"fmt"
	"pulsar-operator/pkg/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeDeployment(c *v1alpha1.PulsarCluster) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeDeploymentName(c),
			Namespace: c.Namespace,
			Labels:    v1alpha1.MakeComponentLabels(c, v1alpha1.BrokerComponent),
		},
		Spec: makeDeploymentSpec(c),
	}
}

func MakeDeploymentName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-broker-deployment", c.GetName())
}

func makeDeploymentSpec(c *v1alpha1.PulsarCluster) appsv1.DeploymentSpec {
	return appsv1.DeploymentSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: v1alpha1.MakeComponentLabels(c, v1alpha1.BrokerComponent),
		},
		Replicas: &c.Spec.Broker.Size,
		Template: makeDeploymentPodTemplate(c),
	}
}

func makeDeploymentPodTemplate(c *v1alpha1.PulsarCluster) v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: c.GetName(),
			Labels:       v1alpha1.MakeComponentLabels(c, v1alpha1.BrokerComponent),
			Annotations:  DeploymentAnnotations,
		},
		Spec: makePodSpec(c),
	}
}
