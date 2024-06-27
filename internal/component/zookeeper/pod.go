package zookeeper

import (
	"fmt"
	"strings"

	"github.com/falser101/pulsar-operator/api/v1alpha1"

	v1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func MakePodDisruptionBudget(c *v1alpha1.PulsarCluster) *policyv1.PodDisruptionBudget {
	count := intstr.FromInt32(1)
	return &policyv1.PodDisruptionBudget{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PodDisruptionBudget",
			APIVersion: "policy/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakePodDisruptionBudgetName(c),
			Namespace: c.Namespace,
		},
		Spec: policyv1.PodDisruptionBudgetSpec{
			MaxUnavailable: &count,
			Selector: &metav1.LabelSelector{
				MatchLabels: v1alpha1.MakeComponentLabels(c, v1alpha1.ZookeeperComponent),
			},
		},
	}
}

func MakePodDisruptionBudgetName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-zookeeper-poddisruptionbudget", c.Name)
}

func makePodSpec(c *v1alpha1.PulsarCluster) v1.PodSpec {
	var p = v1.PodSpec{
		Affinity:           c.Spec.Zookeeper.Pod.Affinity,
		ServiceAccountName: makeServiceAccountName(c),
		SecurityContext:    c.Spec.Zookeeper.Pod.SecurityContext,
		RestartPolicy:      c.Spec.Zookeeper.Pod.RestartPolicy,
		Containers:         []v1.Container{makeContainer(c)},
	}
	if isUseEmptyDirVolume(c) {
		p.Volumes = makeEmptyDirVolume(c)
	}
	return p
}

func makeContainer(c *v1alpha1.PulsarCluster) v1.Container {
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
			{
				Name:      makeZookeeperDateVolumeName(c),
				MountPath: ContainerDataPath,
			},
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
		`bin/apply-config-from-env.py conf/zookeeper.conf;
		bin/apply-config-from-env.py conf/pulsar_env.sh;
		bin/generate-zookeeper-config.sh conf/zookeeper.conf;
		OPTS="${OPTS} -Dlog4j2.formatMsgNoLookups=true" exec bin/pulsar zookeeper;`,
	}
}

func makeContainerPort(c *v1alpha1.PulsarCluster) []v1.ContainerPort {
	return []v1.ContainerPort{
		{
			Name:          "http",
			ContainerPort: c.Spec.Zookeeper.Ports.Http,
			Protocol:      v1.ProtocolTCP,
		},
		{
			Name:          "client",
			ContainerPort: c.Spec.Zookeeper.Ports.Client,
			Protocol:      v1.ProtocolTCP,
		},
		{
			Name:          "follower",
			ContainerPort: c.Spec.Zookeeper.Ports.Follower,
			Protocol:      v1.ProtocolTCP,
		},
		{
			Name:          "leader-election",
			ContainerPort: c.Spec.Zookeeper.Ports.LeaderElection,
			Protocol:      v1.ProtocolTCP,
		},
	}
}

func makeContainerEnv(c *v1alpha1.PulsarCluster) []v1.EnvVar {
	return []v1.EnvVar{
		{
			Name:  ContainerZookeeperServerList,
			Value: strings.Join(makeStatefulSetPodNameList(c), ","),
		},
		{
			Name:  "EXTERNAL_PROVIDED_SERVERS",
			Value: "false",
		},
	}
}

func makeContainerEnvFrom(c *v1alpha1.PulsarCluster) []v1.EnvFromSource {
	return []v1.EnvFromSource{
		{
			ConfigMapRef: &v1.ConfigMapEnvSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: MakeConfigMapName(c),
				},
			},
		},
	}
}

func isUseEmptyDirVolume(c *v1alpha1.PulsarCluster) bool {
	return c.Spec.Zookeeper.Volumes.Data.StorageClassName == ""
}

func MakeWaitZookeeperReadyContainer(c *v1alpha1.PulsarCluster) v1.Container {
	return v1.Container{
		Name:            "wait-zookeeper-ready",
		Image:           c.Spec.Zookeeper.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Zookeeper.Image.PullPolicy,
		Command:         makeContainerCommand(),
		Args:            makeWaitZookeeperReadyContainerCommandArgs(c),
	}
}

func makeWaitZookeeperReadyContainerCommandArgs(c *v1alpha1.PulsarCluster) []string {
	return []string{
		fmt.Sprintf(`until nslookup %s.%s.%s.svc.cluster.local; do
		sleep 3;
		done;`, fmt.Sprintf("%s-%d", MakeStatefulSetName(c), c.Spec.Zookeeper.Replicas-1), MakeServiceName(c), c.Namespace),
	}
}
