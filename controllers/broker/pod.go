package broker

import (
	"k8s.io/api/core/v1"
	cachev1alpha1 "pulsar-operator/api/v1alpha1"
)

func makePodSpec(c *cachev1alpha1.PulsarCluster) v1.PodSpec {
	return v1.PodSpec{
		Affinity:   c.Spec.Broker.Pod.Affinity,
		Containers: []v1.Container{makeContainer(c)},
	}
}

func makeContainer(c *cachev1alpha1.PulsarCluster) v1.Container {
	return v1.Container{
		Name:            "broker",
		Image:           c.Spec.Broker.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Broker.Image.PullPolicy,
		Command:         makeContainerCommand(),
		Args:            makeContainerCommandArgs(),
		Ports:           makeContainerPort(c),
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
		"bin/apply-config-from-env.py conf/broker.conf && " +
			"bin/apply-config-from-env.py conf/pulsar_env.sh && " +
			"bin/gen-yml-from-env.py conf/functions_worker.yml && " +
			"bin/pulsar broker",
	}
}

func makeContainerPort(c *cachev1alpha1.PulsarCluster) []v1.ContainerPort {
	return []v1.ContainerPort{
		{
			Name:          "http",
			ContainerPort: cachev1alpha1.PulsarBrokerHttpServerPort,
			Protocol:      v1.ProtocolTCP,
		},
		{
			Name:          "pulsar",
			ContainerPort: cachev1alpha1.PulsarBrokerPulsarServerPort,
			Protocol:      v1.ProtocolTCP,
		},
	}
}

func makeContainerEnv(c *cachev1alpha1.PulsarCluster) []v1.EnvVar {
	env := make([]v1.EnvVar, 0)
	env = append(env, v1.EnvVar{
		Name:      AdvertisedAddress,
		ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{FieldPath: "status.podIP"}},
	})
	return env
}

func makeContainerEnvFrom(c *cachev1alpha1.PulsarCluster) []v1.EnvFromSource {
	froms := make([]v1.EnvFromSource, 0)

	var configRef v1.ConfigMapEnvSource
	configRef.Name = MakeConfigMapName(c)

	froms = append(froms, v1.EnvFromSource{ConfigMapRef: &configRef})
	return froms
}
