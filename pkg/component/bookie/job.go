package bookie

import (
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"pulsar-operator/pkg/api/v1alpha1"
	"pulsar-operator/pkg/component/zookeeper"
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
	parallelism := int32(1)
	completions := int32(1)
	return batchv1.JobSpec{
		Parallelism: &parallelism,
		Completions: &completions,
		Template:    makePodTemplate(c),
	}
}

func makePodTemplate(c *v1alpha1.PulsarCluster) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: c.GetName(),
			Labels:       v1alpha1.MakeComponentLabels(c, ComponentName),
		},
		Spec: corev1.PodSpec{
			Containers:     []corev1.Container{makeJobContainer(c)},
			InitContainers: []corev1.Container{makeJobInitContainer(c)},
			RestartPolicy:  corev1.RestartPolicyNever,
		},
	}
}

func makeJobInitContainer(c *v1alpha1.PulsarCluster) corev1.Container {
	return corev1.Container{
		Name:            makeJobInitContainerName(c),
		Image:           c.Spec.Bookie.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Bookie.Image.PullPolicy,
		Command:         makeJobInitContainerCommand(),
		Args:            makeJobInitContainerCommandArgs(c),
	}
}

func makeJobInitContainerCommandArgs(c *v1alpha1.PulsarCluster) []string {
	return []string{
		fmt.Sprintf(`until nslookup %s-0.%s.%s.svc.cluster.local; do
            sleep 3;
          done;
            sleep 3;`, zookeeper.MakeStatefulSetName(c), zookeeper.MakeServiceName(c), c.Namespace),
	}
}

func makeJobInitContainerCommand() []string {
	return []string{
		"sh",
		"-c",
	}
}

func makeJobInitContainerName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("wait-zookeeper-ready")
}

func makeJobContainer(c *v1alpha1.PulsarCluster) corev1.Container {
	return corev1.Container{
		Name:            makeJobContainerName(c),
		Image:           c.Spec.Bookie.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Bookie.Image.PullPolicy,
		Command:         makeJobContainerCommand(),
		Args:            makeJobContainerCommandArgs(c),
		EnvFrom:         makeEnvFrom(c),
	}
}

func makeEnvFrom(c *v1alpha1.PulsarCluster) []corev1.EnvFromSource {
	return []corev1.EnvFromSource{
		{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: MakeConfigMapName(c)},
			},
		},
	}
}

func makeJobContainerCommandArgs(c *v1alpha1.PulsarCluster) []string {
	return []string{
		`bin/apply-config-from-env.py conf/bookkeeper.conf;
          if bin/bookkeeper shell whatisinstanceid; then
              echo "bookkeeper cluster already initialized";
          else
              echo "bookkeeper cluster start initialized";
              bin/bookkeeper shell initnewcluster;
          fi`,
	}
}

func makeJobContainerCommand() []string {
	return []string{
		"sh",
		"-c",
	}
}

func makeJobContainerName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-bookie-init", c.Name)
}

func MakeInitBookieJobName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-bookie-init", c.Name)
}
