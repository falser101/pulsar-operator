package bookie

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/internal/component/zookeeper"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeBookieClusterInitJob(c *v1alpha1.PulsarCluster) *batchv1.Job {
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      makeInitBookieClusterJobName(c),
			Namespace: c.Namespace,
			Labels:    v1alpha1.MakeComponentLabels(c, v1alpha1.BookieComponent),
		},
		Spec: makeJobSpec(c),
	}
}

func makeInitBookieClusterJobName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-%s-init", c.Name, v1alpha1.BookieComponent)
}

func makeJobSpec(c *v1alpha1.PulsarCluster) batchv1.JobSpec {
	return batchv1.JobSpec{
		Template: makeBookieInitPodTemplate(c),
	}
}

func makeBookieInitPodTemplate(c *v1alpha1.PulsarCluster) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            fmt.Sprintf("%s-%s-init", c.Name, v1alpha1.BookieComponent),
					Image:           c.Spec.Bookie.Image.GenerateImage(),
					ImagePullPolicy: c.Spec.Bookie.Image.PullPolicy,
					Resources:       corev1.ResourceRequirements{},
					Command: []string{
						"sh",
						"-c",
					},
					Args: []string{
						`bin/apply-config-from-env.py conf/bookkeeper.conf;
						export BOOKIE_MEM="-Xmx128M";
						if bin/bookkeeper shell whatisinstanceid; then
							echo "bookkeeper cluster already initialized";
						else
							bin/bookkeeper shell initnewcluster;
						fi`,
					},
					EnvFrom: []corev1.EnvFromSource{
						{
							ConfigMapRef: &corev1.ConfigMapEnvSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: MakeConfigMapName(c),
								},
							},
						},
					},
				},
			},
			InitContainers: []corev1.Container{
				zookeeper.MakeWaitZookeeperReadyContainer(c),
			},
			RestartPolicy: "OnFailure",
		},
	}
}
