package zookeeper

import (
	"fmt"
	"pulsar-operator/pkg/api/v1alpha1"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
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
			Labels:    v1alpha1.MakeComponentLabels(c, v1alpha1.ZookeeperComponent),
		},
		Spec: makeStatefulSetSpec(c),
	}
}

func MakeStatefulSetName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-zookeeper-statefulset", c.GetName())
}

func makeStatefulSetPodNameList(c *v1alpha1.PulsarCluster) []string {
	result := make([]string, 0)
	for i := 0; i < int(c.Spec.Zookeeper.Size); i++ {
		result = append(result, fmt.Sprintf("%s-%s", MakeStatefulSetName(c), strconv.Itoa(i)))
	}
	return result
}

func makeStatefulSetSpec(c *v1alpha1.PulsarCluster) appsv1.StatefulSetSpec {
	return appsv1.StatefulSetSpec{
		ServiceName: MakeServiceName(c),
		Selector: &metav1.LabelSelector{
			MatchLabels: v1alpha1.MakeComponentLabels(c, v1alpha1.ZookeeperComponent),
		},
		Replicas:            &c.Spec.Zookeeper.Size,
		Template:            makeStatefulSetPodTemplate(c),
		PodManagementPolicy: appsv1.OrderedReadyPodManagement,
		UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
			Type: appsv1.RollingUpdateStatefulSetStrategyType,
		},
	}
}

func makeStatefulSetPodTemplate(c *v1alpha1.PulsarCluster) v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: c.GetName(),
			Labels:       v1alpha1.MakeComponentLabels(c, v1alpha1.ZookeeperComponent),
			Annotations:  StatefulSetAnnotations,
		},
		Spec: makePodSpec(c),
	}
}
