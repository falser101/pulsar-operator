package bookie

import (
	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/internal/component/zookeeper"

	v1 "k8s.io/api/core/v1"
)

func makePodSpec(c *v1alpha1.PulsarCluster) v1.PodSpec {
	p := v1.PodSpec{
		Affinity:   c.Spec.Bookie.Pod.Affinity,
		Containers: []v1.Container{makeContainer(c)},
		InitContainers: []v1.Container{
			zookeeper.MakeWaitZookeeperReadyContainer(c),
			makeInitContainer(c),
		},
	}

	if isUseEmptyDirVolume(c) {
		p.Volumes = makeEmptyDirVolume(c)
	}

	return p
}

func makeContainer(c *v1alpha1.PulsarCluster) v1.Container {
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
			"bin/pulsar bookie",
	}
}

func makeContainerPort(c *v1alpha1.PulsarCluster) []v1.ContainerPort {
	return []v1.ContainerPort{
		{
			Name:          "server",
			ContainerPort: v1alpha1.PulsarBookieServerPort,
			Protocol:      v1.ProtocolTCP,
		},
		{
			Name:          "client",
			ContainerPort: v1alpha1.PulsarBookieClientPort,
			Protocol:      v1.ProtocolTCP,
		},
	}
}

func makeContainerEnv(c *v1alpha1.PulsarCluster) []v1.EnvVar {
	env := make([]v1.EnvVar, 0)
	env = append(env, v1.EnvVar{
		Name: "POD_NAME",
		ValueFrom: &v1.EnvVarSource{
			FieldRef: &v1.ObjectFieldSelector{
				FieldPath:  "metadata.name",
				APIVersion: "v1",
			},
		}})
	env = append(env, v1.EnvVar{
		Name: "POD_NAMESPACE",
		ValueFrom: &v1.EnvVarSource{
			FieldRef: &v1.ObjectFieldSelector{
				FieldPath:  "metadata.namespace",
				APIVersion: "v1",
			},
		}})
	env = append(env, v1.EnvVar{
		Name:  "VOLUME_NAME",
		Value: c.Name + "-bookie-journal",
	})
	env = append(env, v1.EnvVar{
		Name:  "BOOKIE_PORT",
		Value: "3181",
	})
	env = append(env, v1.EnvVar{
		Name:  "BOOKIE_RACK_AWARE_ENABLED",
		Value: "true",
	})
	return env
}

func makeContainerEnvFrom(c *v1alpha1.PulsarCluster) []v1.EnvFromSource {
	froms := make([]v1.EnvFromSource, 0)

	var configRef v1.ConfigMapEnvSource
	configRef.Name = MakeConfigMapName(c)

	froms = append(froms, v1.EnvFromSource{ConfigMapRef: &configRef})
	return froms
}

func makeInitContainer(c *v1alpha1.PulsarCluster) v1.Container {
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
		`set -e; bin/apply-config-from-env.py conf/bookkeeper.conf;until bin/bookkeeper shell whatisinstanceid; do
			sleep 3;
		done;`,
	}
}

func makeInitContainerEnvFrom(c *v1alpha1.PulsarCluster) []v1.EnvFromSource {
	froms := make([]v1.EnvFromSource, 0)

	var configRef v1.ConfigMapEnvSource
	configRef.Name = MakeConfigMapName(c)

	froms = append(froms, v1.EnvFromSource{ConfigMapRef: &configRef})
	return froms
}

func isUseEmptyDirVolume(c *v1alpha1.PulsarCluster) bool {
	return c.Spec.Bookie.StorageClassName == ""
}
