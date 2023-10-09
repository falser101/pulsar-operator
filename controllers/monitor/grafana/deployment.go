package grafana

import (
	"fmt"
	cachev1alpha1 "pulsar-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeDeployment(c *cachev1alpha1.PulsarCluster) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeDeploymentName(c),
			Namespace: c.Namespace,
			Labels:    cachev1alpha1.MakeAllLabels(c, cachev1alpha1.MonitorComponent, cachev1alpha1.MonitorGrafanaComponent),
		},
		Spec: makeDeploymentSpec(c),
	}
}

func MakeDeploymentName(c *cachev1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-monitor-grafana-deployment", c.GetName())
}

func makeDeploymentSpec(c *cachev1alpha1.PulsarCluster) appsv1.DeploymentSpec {
	var replica int32 = 1
	return appsv1.DeploymentSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: cachev1alpha1.MakeAllLabels(c, cachev1alpha1.MonitorComponent, cachev1alpha1.MonitorGrafanaComponent),
		},
		Replicas: &replica,
		Template: makeDeploymentPodTemplate(c),
	}
}

func makeDeploymentPodTemplate(c *cachev1alpha1.PulsarCluster) v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: c.GetName(),
			Labels:       cachev1alpha1.MakeAllLabels(c, cachev1alpha1.MonitorComponent, cachev1alpha1.MonitorGrafanaComponent),
		},
		Spec: makePodSpec(c),
	}
}
