package bookie

import (
	cachev1alpha1 "pulsar-operator/api/v1alpha1"

	"k8s.io/api/core/v1"
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

func makeContainer(c *cachev1alpha1.PulsarCluster) v1.Container {
	return v1.Container{
		Name:            "bookie",
		Image:           c.Spec.Bookie.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Bookie.Image.PullPolicy,
		Command:         makeContainerCommand(),
		Args:            makeContainerCommandArgs(),
		Ports:           makeContainerPort(c),
		Env:             makeContainerEnv(c),
		EnvFrom:         makeContainerEnvFrom(c),

		VolumeMounts: []v1.VolumeMount{
			{
				Name:      makeJournalDataVolumeName(c),
				MountPath: JournalDataMountPath,
			},
			{
				Name:      makeLedgersDataVolumeName(c),
				MountPath: LedgersDataMountPath,
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
		"bin/apply-config-from-env.py conf/bookkeeper.conf && " +
			"bin/apply-config-from-env.py conf/pulsar_env.sh && " +
			"bin/pulsar bookie",
	}
}

func makeContainerPort(c *cachev1alpha1.PulsarCluster) []v1.ContainerPort {
	return []v1.ContainerPort{
		{
			Name:          "client",
			ContainerPort: cachev1alpha1.PulsarBookieServerPort,
			Protocol:      v1.ProtocolTCP,
		},
	}
}

func makeContainerEnv(c *cachev1alpha1.PulsarCluster) []v1.EnvVar {
	env := make([]v1.EnvVar, 0)
	return env
}

func makeContainerEnvFrom(c *cachev1alpha1.PulsarCluster) []v1.EnvFromSource {
	froms := make([]v1.EnvFromSource, 0)

	var configRef v1.ConfigMapEnvSource
	configRef.Name = MakeConfigMapName(c)

	froms = append(froms, v1.EnvFromSource{ConfigMapRef: &configRef})
	return froms
}

func makeInitContainer(c *cachev1alpha1.PulsarCluster) v1.Container {
	return v1.Container{
		Name:            "bookie-metaformat",
		Image:           c.Spec.AutoRecovery.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.AutoRecovery.Image.PullPolicy,
		Command:         makeInitContainerCommand(),
		Args:            makeInitContainerCommandArgs(),
		EnvFrom:         makeInitContainerEnvFrom(c),
	}
}

func makeInitContainerCommand() []string {
	return []string{
		"sh",
		"-c",
	}
}

func makeInitContainerCommandArgs() []string {
	return []string{
		"bin/apply-config-from-env.py conf/bookkeeper.conf && " +
			"bin/bookkeeper shell metaformat --nonInteractive || true;",
	}
}

func makeInitContainerEnvFrom(c *cachev1alpha1.PulsarCluster) []v1.EnvFromSource {
	froms := make([]v1.EnvFromSource, 0)

	var configRef v1.ConfigMapEnvSource
	configRef.Name = MakeConfigMapName(c)

	froms = append(froms, v1.EnvFromSource{ConfigMapRef: &configRef})
	return froms
}

func isUseEmptyDirVolume(c *cachev1alpha1.PulsarCluster) bool {
	return c.Spec.Bookie.StorageClassName == ""
}
