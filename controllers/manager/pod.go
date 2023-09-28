package manager

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/utils/net"
	cachev1alpha1 "pulsar-operator/api/v1alpha1"
	"pulsar-operator/controllers/broker"
)

func makePodSpec(c *cachev1alpha1.PulsarCluster) v1.PodSpec {
	p := v1.PodSpec{
		Affinity:       c.Spec.Bookie.Pod.Affinity,
		Containers:     []v1.Container{makeContainer(c)},
		InitContainers: []v1.Container{makeInitContainer(c)},
	}

	if isUseEmptyDirVolume(c) {
		p.Volumes = makeEmptyDirVolume(c)
	}

	return p
}

func isUseEmptyDirVolume(c *cachev1alpha1.PulsarCluster) bool {
	return c.Spec.Manager.StorageClassName == ""
}

func makeInitContainer(c *cachev1alpha1.PulsarCluster) v1.Container {
	return v1.Container{
		Name:            "wait-broker-ready",
		Image:           c.Spec.AutoRecovery.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.AutoRecovery.Image.PullPolicy,
		Command:         makeInitContainerCommand(),
		Args:            makeInitContainerCommandArgs(c),
	}
}

func makeInitContainerCommandArgs(c *cachev1alpha1.PulsarCluster) []string {
	return []string{
		fmt.Sprintf(" brokerServiceNumber=\"$(nslookup -timeout=10 %s | grep Name | wc -l)\"; until [ ${brokerServiceNumber} -ge 1 ]; do\n"+
			"            echo \"broker cluster %s isn't ready yet ... check in 10 seconds ...\";\n"+
			"            sleep 10;\n"+
			"            brokerServiceNumber=\"$(nslookup -timeout=10 %s | grep Name | wc -l)\";\n"+
			"          done; echo \"broker cluster is ready\"", broker.MakeServiceName(c), c.Name, broker.MakeServiceName(c)),
	}

}

func makeInitContainerCommand() []string {
	return []string{
		"sh",
		"-c",
	}
}

func makeContainer(c *cachev1alpha1.PulsarCluster) v1.Container {
	return v1.Container{
		Name:            "manager",
		Image:           c.Spec.Manager.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Manager.Image.PullPolicy,
		Command:         makeContainerCommand(),
		Args:            makeContainerCommandArgs(),
		Ports:           makeContainerPort(c),
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      makeManagerDataName(c),
				MountPath: "/data",
			},
			{
				Name:      makeManagerScriptName(c),
				MountPath: "/pulsar-manager/pulsar-manager.sh",
				SubPath:   EntrypointKey,
			},
			{
				Name:      makeManagerBackendScriptName(c),
				MountPath: "/pulsar-manager/pulsar-backend-entrypoint.sh",
				SubPath:   BackendEntrypointKey,
			},
			{
				Name:      makeManagerTokenKeysName(c),
				MountPath: "/pulsar/keys",
				ReadOnly:  true,
			},
			{
				Name:      makeManagerTokensName(c),
				MountPath: "/pulsar/tokens",
				ReadOnly:  true,
			},
		},
	}
}

func makeContainerPort(c *cachev1alpha1.PulsarCluster) []v1.ContainerPort {
	return []v1.ContainerPort{
		{
			ContainerPort: cachev1alpha1.PulsarManagerFrontendPort,
			Name:          "frontend",
			Protocol:      v1.Protocol(net.TCP),
		},
		{
			ContainerPort: cachev1alpha1.PulsarManagerBackendPort,
			Name:          "backend",
			Protocol:      v1.Protocol(net.TCP),
		},
	}
}

func makeContainerCommandArgs() []string {
	return []string{
		"/pulsar-manager/pulsar-manager.sh",
	}
}

func makeContainerCommand() []string {
	return []string{
		"sh",
		"-c",
	}
}
