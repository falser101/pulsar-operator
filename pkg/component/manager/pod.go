package manager

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/pkg/component/broker"
	v1 "k8s.io/api/core/v1"
	"k8s.io/utils/net"
)

func makePodSpec(c *v1alpha1.Pulsar) v1.PodSpec {
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

func isUseEmptyDirVolume(c *v1alpha1.Pulsar) bool {
	return c.Spec.Manager.StorageClassName == ""
}

func makeInitContainer(c *v1alpha1.Pulsar) v1.Container {
	return v1.Container{
		Name:            "wait-broker-ready",
		Image:           c.Spec.AutoRecovery.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.AutoRecovery.Image.PullPolicy,
		Command:         makeInitContainerCommand(),
		Args:            makeInitContainerCommandArgs(c),
	}
}

func makeInitContainerCommandArgs(c *v1alpha1.Pulsar) []string {
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

func makeContainer(c *v1alpha1.Pulsar) v1.Container {
	var volumeMounts = []v1.VolumeMount{
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
	}
	if c.Spec.Broker.Authentication.Enabled {
		volumeMounts = append(volumeMounts,
			v1.VolumeMount{
				Name:      makeManagerTokenKeysName(c),
				MountPath: "/pulsar/keys",
				ReadOnly:  true,
			},
			v1.VolumeMount{
				Name:      makeManagerTokensName(c),
				MountPath: "/pulsar/tokens",
				ReadOnly:  true,
			})
	}
	return v1.Container{
		Name:            "manager",
		Image:           c.Spec.Manager.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Manager.Image.PullPolicy,
		Command:         makeContainerCommand(),
		Args:            makeContainerCommandArgs(),
		Ports:           makeContainerPort(c),
		VolumeMounts:    volumeMounts,
	}
}

func makeContainerPort(c *v1alpha1.Pulsar) []v1.ContainerPort {
	return []v1.ContainerPort{
		{
			ContainerPort: v1alpha1.PulsarManagerFrontendPort,
			Name:          "frontend",
			Protocol:      v1.Protocol(net.TCP),
		},
		{
			ContainerPort: v1alpha1.PulsarManagerBackendPort,
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
