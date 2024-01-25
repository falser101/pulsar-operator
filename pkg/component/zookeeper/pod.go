package zookeeper

import (
	"fmt"
	"strings"

	"github.com/falser101/pulsar-operator/api/v1alpha1"

	v1 "k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func MakePodDisruptionBudget(c *v1alpha1.Pulsar) *v1beta1.PodDisruptionBudget {
	count := intstr.FromInt(1)
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
				MatchLabels: v1alpha1.MakeComponentLabels(c, v1alpha1.ZookeeperComponent),
			},
		},
	}
}

func MakePodDisruptionBudgetName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-zookeeper-poddisruptionbudget", c.GetName())
}

func makePodSpec(c *v1alpha1.Pulsar) v1.PodSpec {
	var p = v1.PodSpec{
		Affinity:   c.Spec.Zookeeper.Pod.Affinity,
		Containers: []v1.Container{makeContainer(c)},
	}
	if isUseEmptyDirVolume(c) {
		p.Volumes = makeEmptyDirVolume(c)
	}
	return p
}

func makeContainer(c *v1alpha1.Pulsar) v1.Container {
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
			{Name: makeZookeeperDateVolumeName(c), MountPath: ContainerDataPath},
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

func makeContainerPort(c *v1alpha1.Pulsar) []v1.ContainerPort {
	return []v1.ContainerPort{
		{
			Name:          "client",
			ContainerPort: v1alpha1.ZookeeperContainerClientDefaultPort,
			Protocol:      v1.ProtocolTCP,
		},
		{
			Name:          "server",
			ContainerPort: v1alpha1.ZookeeperContainerServerDefaultPort,
			Protocol:      v1.ProtocolTCP,
		},
		{
			Name:          "leader-election",
			ContainerPort: v1alpha1.ZookeeperContainerLeaderElectionPort,
			Protocol:      v1.ProtocolTCP,
		},
	}
}

func makeContainerEnv(c *v1alpha1.Pulsar) []v1.EnvVar {
	env := make([]v1.EnvVar, 0)
	env = append(env, v1.EnvVar{Name: ContainerZookeeperServerList, Value: strings.Join(makeStatefulSetPodNameList(c), ",")})
	return env
}

func makeContainerEnvFrom(c *v1alpha1.Pulsar) []v1.EnvFromSource {
	froms := make([]v1.EnvFromSource, 0)

	var configRef v1.ConfigMapEnvSource
	configRef.Name = MakeConfigMapName(c)

	froms = append(froms, v1.EnvFromSource{ConfigMapRef: &configRef})
	return froms
}

func isUseEmptyDirVolume(c *v1alpha1.Pulsar) bool {
	return c.Spec.Zookeeper.StorageClassName == ""
}

func MakeWaitZookeeperReadyContainer(c *v1alpha1.Pulsar) v1.Container {
	return v1.Container{
		Name:            "wait-zookeeper-ready",
		Image:           c.Spec.Zookeeper.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Zookeeper.Image.PullPolicy,
		Command:         makeContainerCommand(),
		Args:            makeWaitZookeeperReadyContainerCommandArgs(c),
		EnvFrom:         makeWaitZookeeperReadyContainerEnvFrom(c),
	}
}

func makeWaitZookeeperReadyContainerCommandArgs(c *v1alpha1.Pulsar) []string {
	return []string{
		fmt.Sprintf(`until bin/pulsar zookeeper-shell -server %s-zookeeper-service get /admin/clusters/%s; do
		sleep 3;
		done;`, c.GetName(), c.GetName()),
	}
}

func makeWaitZookeeperReadyContainerEnvFrom(c *v1alpha1.Pulsar) []v1.EnvFromSource {
	return []v1.EnvFromSource{
		{ConfigMapRef: &v1.ConfigMapEnvSource{LocalObjectReference: v1.LocalObjectReference{Name: MakeConfigMapName(c)}}},
	}
}
