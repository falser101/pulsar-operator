package autorecovery

import (
	"k8s.io/api/core/v1"
	cachev1alpha1 "pulsar-operator/api/v1alpha1"
	"pulsar-operator/controllers/bookie"
)

func makePodSpec(c *cachev1alpha1.PulsarCluster) v1.PodSpec {
	return v1.PodSpec{
		Containers: []v1.Container{makeContainer(c)},
	}
}

func makeContainer(c *cachev1alpha1.PulsarCluster) v1.Container {
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

func makeContainerEnv(c *cachev1alpha1.PulsarCluster) []v1.EnvVar {
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

func makeContainerEnvFrom(c *cachev1alpha1.PulsarCluster) []v1.EnvFromSource {
	froms := make([]v1.EnvFromSource, 0)

	var configRef v1.ConfigMapEnvSource
	configRef.Name = bookie.MakeConfigMapName(c)

	froms = append(froms, v1.EnvFromSource{ConfigMapRef: &configRef})
	return froms
}
