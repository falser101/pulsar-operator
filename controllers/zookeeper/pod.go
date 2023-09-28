package zookeeper

import (
	"fmt"
	cachev1alpha1 "pulsar-operator/api/v1alpha1"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func MakePodDisruptionBudget(c *cachev1alpha1.PulsarCluster) *v1beta1.PodDisruptionBudget {
	count := intstr.FromInt32(1)
	return &v1beta1.PodDisruptionBudget{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PodDisruptionBudget",
			APIVersion: "policy/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakePodDisruptionBudgetName(c),
			Namespace: c.Namespace,
		},
		Spec: v1beta1.PodDisruptionBudgetSpec{
			MaxUnavailable: &count,
			Selector: &metav1.LabelSelector{
				MatchLabels: cachev1alpha1.MakeComponentLabels(c, cachev1alpha1.ZookeeperComponent),
			},
		},
	}
}

func MakePodDisruptionBudgetName(c *cachev1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-zookeeper-poddisruptionbudget", c.GetName())
}

func makePodSpec(c *cachev1alpha1.PulsarCluster) v1.PodSpec {
	return v1.PodSpec{
		Affinity:   c.Spec.Zookeeper.Pod.Affinity,
		Containers: []v1.Container{makeContainer(c)},
		Volumes: []v1.Volume{
			{
				Name: ContainerDataVolumeName,
				VolumeSource: v1.VolumeSource{
					EmptyDir: &v1.EmptyDirVolumeSource{},
				},
			},
		},
	}
}

func makeContainer(c *cachev1alpha1.PulsarCluster) v1.Container {
	return v1.Container{
		Name:            "zookeeper",
		Image:           c.Spec.Zookeeper.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Zookeeper.Image.PullPolicy,
		Command:         makeContainerCommand(),
		Args:            makeContainerCommandArgs(),
		Ports:           makeContainerPort(c),
		Env:             makeContainerEnv(c),
		EnvFrom:         makeContainerEnvFrom(c),

		ReadinessProbe: &v1.Probe{
			InitialDelaySeconds: 5,
			TimeoutSeconds:      5,
			ProbeHandler:        v1.ProbeHandler{Exec: &v1.ExecAction{Command: []string{ReadinessProbeScript}}},
		},
		LivenessProbe: &v1.Probe{
			InitialDelaySeconds: 15,
			TimeoutSeconds:      5,
			ProbeHandler: v1.ProbeHandler{
				Exec: &v1.ExecAction{Command: []string{LivenessProbeScript}},
			},
		},

		VolumeMounts: []v1.VolumeMount{
			{Name: ContainerDataVolumeName, MountPath: ContainerDataPath},
		},
	}
}

func makeContainerCommand() []string {
	return []string{
		"sh",
		"-c",
	}
}

func makeContainerCommandArgs() []string {
	return []string{
		"bin/apply-config-from-env.py conf/zookeeper.conf && " +
			"bin/apply-config-from-env.py conf/pulsar_env.sh && " +
			"bin/generate-zookeeper-config.sh conf/zookeeper.conf && " +
			"bin/pulsar zookeeper",
	}
}

func makeContainerPort(c *cachev1alpha1.PulsarCluster) []v1.ContainerPort {
	return []v1.ContainerPort{
		{
			Name:          "client",
			ContainerPort: cachev1alpha1.ZookeeperContainerClientDefaultPort,
			Protocol:      v1.ProtocolTCP,
		},
		{
			Name:          "server",
			ContainerPort: cachev1alpha1.ZookeeperContainerServerDefaultPort,
			Protocol:      v1.ProtocolTCP,
		},
		{
			Name:          "leader-election",
			ContainerPort: cachev1alpha1.ZookeeperContainerLeaderElectionPort,
			Protocol:      v1.ProtocolTCP,
		},
	}
}

func makeContainerEnv(c *cachev1alpha1.PulsarCluster) []v1.EnvVar {
	env := make([]v1.EnvVar, 0)
	env = append(env, v1.EnvVar{Name: ContainerZookeeperServerList, Value: strings.Join(makeStatefulSetPodNameList(c), ",")})
	return env
}

func makeContainerEnvFrom(c *cachev1alpha1.PulsarCluster) []v1.EnvFromSource {
	froms := make([]v1.EnvFromSource, 0)

	var configRef v1.ConfigMapEnvSource
	configRef.Name = MakeConfigMapName(c)

	froms = append(froms, v1.EnvFromSource{ConfigMapRef: &configRef})
	return froms
}
