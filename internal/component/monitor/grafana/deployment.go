package grafana

import (
	"fmt"
	"github.com/falser101/pulsar-operator/api/v1alpha1"

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
			Labels:    v1alpha1.MakeAllLabels(c, v1alpha1.MonitorComponent, v1alpha1.MonitorGrafanaComponent),
		},
		Spec: makeDeploymentSpec(c),
	}
}

func MakeDeploymentName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-monitor-grafana-deployment", c.GetName())
}

func makeDeploymentSpec(c *v1alpha1.PulsarCluster) appsv1.DeploymentSpec {
	var replica int32 = 1
	return appsv1.DeploymentSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: v1alpha1.MakeAllLabels(c, v1alpha1.MonitorComponent, v1alpha1.MonitorGrafanaComponent),
		},
		Replicas: &replica,
		Template: makeDeploymentPodTemplate(c),
	}
}

func makeDeploymentPodTemplate(c *v1alpha1.PulsarCluster) v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: c.GetName(),
			Labels:       v1alpha1.MakeAllLabels(c, v1alpha1.MonitorComponent, v1alpha1.MonitorGrafanaComponent),
		},
		Spec: makePodSpec(c),
	}
}
