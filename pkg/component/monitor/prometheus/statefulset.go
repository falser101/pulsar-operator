package prometheus

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"pulsar-operator/pkg/api/v1alpha1"
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
			Labels:    v1alpha1.MakeAllLabels(c, v1alpha1.MonitorComponent, v1alpha1.MonitorPrometheusComponent),
		},
		Spec: makeStatefulSetSpec(c),
	}
}

func makeStatefulSetSpec(c *v1alpha1.PulsarCluster) appsv1.StatefulSetSpec {
	spec := appsv1.StatefulSetSpec{
		ServiceName: MakeServiceName(c),
		Selector: &metav1.LabelSelector{
			MatchLabels: v1alpha1.MakeAllLabels(c, v1alpha1.MonitorComponent, v1alpha1.MonitorPrometheusComponent),
		},
		Replicas:            &c.Spec.Monitor.Prometheus.Size,
		Template:            makeStatefulSetPodTemplate(c),
		PodManagementPolicy: appsv1.OrderedReadyPodManagement,
		UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
			Type: appsv1.RollingUpdateStatefulSetStrategyType,
		},
	}
	if !isUseEmptyDirVolume(c) {
		spec.VolumeClaimTemplates = makeVolumeClaimTemplates(c)
	}
	return spec
}

func makeStatefulSetPodTemplate(c *v1alpha1.PulsarCluster) v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: c.GetName(),
			Labels:       v1alpha1.MakeAllLabels(c, v1alpha1.MonitorComponent, v1alpha1.MonitorPrometheusComponent),
		},
		Spec: makePodSpec(c),
	}
}

func MakeStatefulSetName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-monitor-prometheus-statefulset", c.Name)
}
