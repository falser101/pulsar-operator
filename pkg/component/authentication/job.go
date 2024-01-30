package authentication

import (
	"github.com/falser101/pulsar-operator/api/v1alpha1"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeJob(c *v1alpha1.Pulsar) *v1.Job {
	var backoffLimit int32 = 6
	return &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Name + "-authentication-job",
			Namespace: c.Namespace,
		},
		Spec: v1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					ServiceAccountName: MakeServiceAccountName(c),
					Containers:         []corev1.Container{makeContainer(c)},
					RestartPolicy:      corev1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backoffLimit,
		},
	}
}

func makeContainer(c *v1alpha1.Pulsar) corev1.Container {
	return corev1.Container{
		Name:    "prepare-secrets",
		Image:   "",
		Command: []string{"sh", "-c"},
		Args:    []string{"/pulsar/scripts/pulsar/prepare_helm_release.sh -n default -k pulsar"},
	}
}
