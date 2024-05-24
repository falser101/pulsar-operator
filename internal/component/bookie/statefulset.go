package bookie

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
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeStatefulSetName(c),
			Namespace: c.Namespace,
			Labels:    v1alpha1.MakeComponentLabels(c, v1alpha1.BookieComponent),
		},
		Spec: makeStatefulSetSpec(c),
	}
}

func MakeStatefulSetName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-%s", c.Name, v1alpha1.BookieComponent)
}

func makeStatefulSetSpec(c *v1alpha1.PulsarCluster) appsv1.StatefulSetSpec {
	s := appsv1.StatefulSetSpec{
		ServiceName: MakeServiceName(c),
		Selector: &metav1.LabelSelector{
			MatchLabels: v1alpha1.MakeComponentLabels(c, v1alpha1.BookieComponent),
		},
		Replicas:            &c.Spec.Bookie.Replicas,
		Template:            makeStatefulSetPodTemplate(c),
		PodManagementPolicy: appsv1.OrderedReadyPodManagement,
		UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
			Type: appsv1.RollingUpdateStatefulSetStrategyType,
		},
	}

	if !isUseEmptyDirVolume(c) {
		s.VolumeClaimTemplates = makeVolumeClaimTemplates(c)
	}

	return s
}

func makeStatefulSetPodTemplate(c *v1alpha1.PulsarCluster) v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: c.Name,
			Labels:       v1alpha1.MakeComponentLabels(c, v1alpha1.BookieComponent),
			Annotations:  StatefulSetAnnotations,
		},
		Spec: makePodSpec(c),
	}
}
