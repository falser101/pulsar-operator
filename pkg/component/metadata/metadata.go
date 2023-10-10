package metadata

import (
	"fmt"
	"k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1alpha12 "pulsar-operator/pkg/api/v1alpha1"
	"pulsar-operator/pkg/component/broker"
	"pulsar-operator/pkg/component/zookeeper"
)

const (
	JobContainerName = "pulsar-cluster-metadata-init-container"

	ComponentName = "init-cluster-metadata-job"
)

func MakeInitClusterMetaDataJob(c *v1alpha12.PulsarCluster) *v1.Job {
	return &v1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeInitClusterMetaDataJobName(c),
			Namespace: c.Namespace,
			Labels:    v1alpha12.MakeComponentLabels(c, ComponentName),
		},
		Spec: makeJobSpec(c),
	}
}

func MakeInitClusterMetaDataJobName(c *v1alpha12.PulsarCluster) string {
	return fmt.Sprintf("%s-init-cluster-metadata-job", c.GetName())
}

func makeJobSpec(c *v1alpha12.PulsarCluster) v1.JobSpec {
	return v1.JobSpec{
		Template: makePodTemplate(c),
	}
}

func makePodTemplate(c *v1alpha12.PulsarCluster) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: c.GetName(),
			Labels:       v1alpha12.MakeComponentLabels(c, ComponentName),
		},
		Spec: corev1.PodSpec{
			Containers:    []corev1.Container{makeContainer(c)},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
}

func makeContainer(c *v1alpha12.PulsarCluster) corev1.Container {
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

func makeContainerCommandArgs(c *v1alpha12.PulsarCluster) []string {
	brokerServiceName := broker.MakeServiceName(c)
	webServiceUrl := fmt.Sprintf("http://%s.%s.%s:%d",
		brokerServiceName, c.Namespace, v1alpha12.ServiceDomain, v1alpha12.PulsarBrokerHttpServerPort)
	brokerServiceUrl := fmt.Sprintf("pulsar://%s.%s.%s:%d",
		brokerServiceName, c.Namespace, v1alpha12.ServiceDomain, v1alpha12.PulsarBrokerPulsarServerPort)
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
