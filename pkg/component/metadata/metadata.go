package metadata

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/pkg/component/broker"
	"github.com/falser101/pulsar-operator/pkg/component/zookeeper"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	JobContainerName     = "pulsar-cluster-metadata-job"
	JobInitContainerName = JobContainerName + "-init"
	ComponentName        = "init-cluster-metadata-job"
)

func MakeInitClusterMetaDataJob(c *v1alpha1.Pulsar) *v1.Job {
	return &v1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeInitClusterMetaDataJobName(c),
			Namespace: c.Namespace,
			Labels:    v1alpha1.MakeComponentLabels(c, ComponentName),
		},
		Spec: makeJobSpec(c),
	}
}

func MakeInitClusterMetaDataJobName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-init-cluster-metadata-job", c.GetName())
}

func makeJobSpec(c *v1alpha1.Pulsar) v1.JobSpec {
	return v1.JobSpec{
		Template: makePodTemplate(c),
	}
}

func makePodTemplate(c *v1alpha1.Pulsar) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: c.GetName(),
			Labels:       v1alpha1.MakeComponentLabels(c, ComponentName),
		},
		Spec: corev1.PodSpec{
			InitContainers: []corev1.Container{makeInitContainer(c)},
			Containers:     []corev1.Container{makeContainer(c)},
			RestartPolicy:  corev1.RestartPolicyNever,
		},
	}
}

func makeInitContainer(c *v1alpha1.Pulsar) corev1.Container {
	return corev1.Container{
		Name:            JobInitContainerName,
		Image:           c.Spec.Zookeeper.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Zookeeper.Image.PullPolicy,
		Command:         makeContainerCommand(),
		Args:            makeContainerCommandArgs(c),
	}
}

func makeContainer(c *v1alpha1.Pulsar) corev1.Container {
	return corev1.Container{
		Name:            JobContainerName,
		Image:           c.Spec.Zookeeper.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Zookeeper.Image.PullPolicy,
		Command:         makeContainerCommand(),
		Args:            makeContainerCommandArgs(c),
	}
}

func makeContainerCommand() []string {
	return []string{
		"sh",
		"-c",
	}
}

func makeContainerCommandArgs(c *v1alpha1.Pulsar) []string {
	brokerServiceName := broker.MakeServiceName(c)
	webServiceUrl := fmt.Sprintf("http://%s.%s.%s:%d",
		brokerServiceName, c.Namespace, v1alpha1.ServiceDomain, v1alpha1.PulsarBrokerHttpServerPort)
	brokerServiceUrl := fmt.Sprintf("pulsar://%s.%s.%s:%d",
		brokerServiceName, c.Namespace, v1alpha1.ServiceDomain, v1alpha1.PulsarBrokerPulsarServerPort)
	return []string{
		"bin/pulsar initialize-cluster-metadata " +
			fmt.Sprintf("--cluster %s ", c.GetName()) +
			fmt.Sprintf("--zookeeper %s ", zookeeper.MakeServiceName(c)) +
			fmt.Sprintf("--configuration-store %s ", zookeeper.MakeServiceName(c)) +
			fmt.Sprintf(" --web-service-url %s/ ", webServiceUrl) +
			fmt.Sprintf("--broker-service-url %s/ ", brokerServiceUrl) +
			"|| true;",
	}
}
