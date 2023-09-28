package manager

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cachev1alpha1 "pulsar-operator/api/v1alpha1"
)

func MakeStatefulSet(c *cachev1alpha1.PulsarCluster) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeStatefulSetName(c),
			Namespace: c.Namespace,
			Labels:    cachev1alpha1.MakeComponentLabels(c, cachev1alpha1.ManagerComponent),
		},
		Spec: makeStatefulSetSpec(c),
	}
}

func makeStatefulSetSpec(c *cachev1alpha1.PulsarCluster) appsv1.StatefulSetSpec {
	s := appsv1.StatefulSetSpec{
		ServiceName: MakeServiceName(c),
		Selector: &metav1.LabelSelector{
			MatchLabels: cachev1alpha1.MakeComponentLabels(c, cachev1alpha1.ManagerComponent),
		},
		Replicas:            &c.Spec.Manager.Size,
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

func makeStatefulSetPodTemplate(c *cachev1alpha1.PulsarCluster) v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: c.GetName(),
			Labels:       cachev1alpha1.MakeComponentLabels(c, cachev1alpha1.ManagerComponent),
		},
		Spec: makePodSpec(c),
	}
}

func MakeStatefulSetName(c *cachev1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-manager-statefulset", c.GetName())
}
