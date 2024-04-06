package autorecovery

import (
	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/internal/component/bookie"
	"k8s.io/api/core/v1"
)

func makePodSpec(c *v1alpha1.PulsarCluster) v1.PodSpec {
	return v1.PodSpec{
		Containers: []v1.Container{makeContainer(c)},
	}
}

func makeContainer(c *v1alpha1.PulsarCluster) v1.Container {
	return v1.Container{
		Name:            "replication-worker",
		Image:           c.Spec.AutoRecovery.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.AutoRecovery.Image.PullPolicy,
		Command:         makeContainerCommand(),
		Args:            makeContainerCommandArgs(),
		Env:             makeContainerEnv(c),
		EnvFrom:         makeContainerEnvFrom(c),
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
			"bin/bookkeeper autorecovery",
	}
}

func makeContainerEnv(c *v1alpha1.PulsarCluster) []v1.EnvVar {
	env := make([]v1.EnvVar, 0)
	env = append(env,
		v1.EnvVar{
			Name:  "BOOKIE_MEM",
			Value: BookieMemData,
		},
		v1.EnvVar{
			Name:  "PULSAR_GC",
			Value: PulsarGCData,
		},
	)
	return env
}

func makeContainerEnvFrom(c *v1alpha1.PulsarCluster) []v1.EnvFromSource {
	froms := make([]v1.EnvFromSource, 0)

	var configRef v1.ConfigMapEnvSource
	configRef.Name = bookie.MakeConfigMapName(c)

	froms = append(froms, v1.EnvFromSource{ConfigMapRef: &configRef})
	return froms
}
