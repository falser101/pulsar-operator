package zookeeper

import (
	"fmt"
	"strconv"

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
			Name:      MakeName(c),
			Namespace: c.Namespace,
			Labels:    v1alpha1.MakeComponentLabels(c, v1alpha1.ZookeeperComponent),
		},
		Spec: makeStatefulSetSpec(c),
	}
}

func makeStatefulSetPodNameList(c *v1alpha1.PulsarCluster) []string {
	result := make([]string, 0)
	for i := 0; i < int(c.Spec.Zookeeper.Replicas); i++ {
		result = append(result, fmt.Sprintf("%s-%s", MakeName(c), strconv.Itoa(i)))
	}
	return result
}

func makeStatefulSetSpec(c *v1alpha1.PulsarCluster) appsv1.StatefulSetSpec {
	var spec = appsv1.StatefulSetSpec{
		ServiceName: MakeName(c),
		Replicas:    &c.Spec.Zookeeper.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: v1alpha1.MakeComponentLabels(c, v1alpha1.ZookeeperComponent),
		},
		UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
			Type: appsv1.RollingUpdateStatefulSetStrategyType,
		},
		PodManagementPolicy: appsv1.ParallelPodManagement,
		Template:            makeStatefulSetPodTemplate(c),
	}
	if !isUseEmptyDirVolume(c) {
		spec.VolumeClaimTemplates = makeVolumeClaimTemplates(c)
	}
	return spec
}

func makeStatefulSetPodTemplate(c *v1alpha1.PulsarCluster) v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      v1alpha1.MakeComponentLabels(c, v1alpha1.ZookeeperComponent),
			Annotations: StatefulSetAnnotations,
		},
		Spec: makePodSpec(c),
	}
}
