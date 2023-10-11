package bookie

import (
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"pulsar-operator/pkg/api/v1alpha1"
)

func MakeJob(c *v1alpha1.PulsarCluster) *batchv1.Job {
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeInitBookieJobName(c),
			Namespace: c.Namespace,
			Labels:    v1alpha1.MakeComponentLabels(c, ComponentName),
		},
		Spec: makeJobSpec(c),
	}
}

func makeJobSpec(c *v1alpha1.PulsarCluster) batchv1.JobSpec {
	return batchv1.JobSpec{}
}

func MakeInitBookieJobName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-bookie-init", c.Name)
}
