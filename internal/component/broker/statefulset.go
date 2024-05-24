package broker

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeStatefulSet(c *v1alpha1.PulsarCluster) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      makeStatefulSetName(c),
			Namespace: c.Namespace,
			Labels:    v1alpha1.MakeComponentLabels(c, v1alpha1.BrokerComponent),
		},
		Spec: makeStatefulSetSpec(c),
	}
}

func makeStatefulSetName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-%s", c.Name, v1alpha1.BrokerComponent)
}

func makeStatefulSetSpec(c *v1alpha1.PulsarCluster) appsv1.StatefulSetSpec {
	return appsv1.StatefulSetSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: v1alpha1.MakeComponentLabels(c, v1alpha1.BrokerComponent),
		},
		Replicas: &c.Spec.Broker.Replicas,
		Template: makeStatefulSetPodTemplate(c),
	}
}

func makeStatefulSetPodTemplate(c *v1alpha1.PulsarCluster) v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: c.Name,
			Labels:       v1alpha1.MakeComponentLabels(c, v1alpha1.BrokerComponent),
			Annotations:  StatefulSetAnnotations,
		},
		Spec: makePodSpec(c),
	}
}
