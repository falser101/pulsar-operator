package manager

import (
	"fmt"
	"github.com/falser101/pulsar-operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ComponentName = "init-manager-job"
)

func MakeJob(c *v1alpha1.Pulsar) *batchv1.Job {
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeInitManagerJobName(c),
			Namespace: c.Namespace,
			Labels:    v1alpha1.MakeComponentLabels(c, ComponentName),
		},
		Spec: makeJobSpec(c),
	}
}

func MakeInitManagerJobName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-init-manager-job", c.GetName())
}

func makeJobSpec(c *v1alpha1.Pulsar) batchv1.JobSpec {
	parallelism := int32(1)
	completions := int32(1)
	return batchv1.JobSpec{
		Parallelism: &parallelism,
		Completions: &completions,
		Template:    makePodTemplate(c),
	}
}

func makePodTemplate(c *v1alpha1.Pulsar) corev1.PodTemplateSpec {
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

func makeJobInitContainer(c *v1alpha1.Pulsar) corev1.Container {
	return corev1.Container{
		Name:            makeJobInitContainerName(c),
		Image:           c.Spec.Broker.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Broker.Image.PullPolicy,
		Command:         makeJobInitContainerCommand(),
		Args:            makeJobInitContainerCommandArgs(c),
	}
}

func makeJobInitContainerCommandArgs(c *v1alpha1.Pulsar) []string {
	return []string{
		fmt.Sprintf("pmServiceNumber=\"$(nslookup -timeout=10 %s | grep Name | wc -l)\"; until [ ${pmServiceNumber} -ge 1 ]; do\n"+
			"            echo \"Pulsar Manager cluster %s isn't ready yet ... check in 10 seconds ...\";\n"+
			"            sleep 10;\n"+
			"            pmServiceNumber=\"$(nslookup -timeout=10 %s | grep Name | wc -l)\";\n"+
			"          done; sleep 5; echo \"Pulsar Manager cluster is ready\";", MakeServiceName(c), c.Name, MakeServiceName(c)),
	}
}

func makeJobInitContainerCommand() []string {
	return []string{
		"sh",
		"-c",
	}
}

func makeJobInitContainerName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("wait-pulsar-manager-ready")
}

func makeJobContainer(c *v1alpha1.Pulsar) corev1.Container {
	return corev1.Container{
		Name:            makeJobContainerName(c),
		Image:           c.Spec.Manager.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Manager.Image.PullPolicy,
		Command:         makeJobContainerCommand(),
		Args:            makeJobContainerCommandArgs(c),
	}
}

func makeJobContainerName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-init-manager-container", c.Name)
}

func makeJobContainerCommandArgs(c *v1alpha1.Pulsar) []string {
	return []string{
		fmt.Sprintf("apk add curl; export CSRF_TOKEN=$(curl http://%s:7750/pulsar-manager/csrf-token);"+
			"          curl -H \"Content-Type: application/json\" \\\n     -H \"X-XSRF-TOKEN: $CSRF_TOKEN\""+
			"          \\\n     -H \"Cookie: XSRF-TOKEN=$CSRF_TOKEN;\" \\\n     -X PUT \\\n     http://%s:7750/pulsar-manager/users/superuser"+
			"          \\\n     -d '{\"name\": \"tonglinkq\", \"password\": \"tonglinkq@123\","+
			"          \"description\": \"Pulsar Manager Admin\", \"email\": \"support@pulsar.io\"}'"+
			"          \n", MakeServiceName(c), MakeServiceName(c)),
	}
}

func makeJobContainerCommand() []string {
	return []string{
		"sh",
		"-c",
	}
}
