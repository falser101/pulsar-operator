package broker

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/pkg/component/bookie"
	v1 "k8s.io/api/core/v1"
)

func makePodSpec(c *v1alpha1.Pulsar) v1.PodSpec {
	return v1.PodSpec{
		Affinity:       c.Spec.Broker.Pod.Affinity,
		InitContainers: []v1.Container{makeWaitBookieReadyContainer(c)},
		Containers:     []v1.Container{makeContainer(c)},
	}
}

func makeWaitBookieReadyContainer(c *v1alpha1.Pulsar) v1.Container {
	return v1.Container{
		Name:            "wait-bookie-ready",
		Image:           c.Spec.Broker.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Broker.Image.PullPolicy,
		Command:         []string{"sh", "-c"},
		Args: []string{
			fmt.Sprintf(" response=\"$(curl -s %s:8000/heartbeat)\";\n"+
				"        until [ \"$response\" = \"OK\" ]; do\n"+
				"            echo \"$response, bookie isn't ready\";\n"+
				"            sleep 1;\n"+
				"            response=\"$(curl -s %s:8000/heartbeat)\";\n"+
				"        done; echo \"$response, bookie is ready\"", bookie.MakeServiceName(c), bookie.MakeServiceName(c)),
		}}
}

func makeContainer(c *v1alpha1.Pulsar) v1.Container {
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

func makeContainerPort(c *v1alpha1.Pulsar) []v1.ContainerPort {
	return []v1.ContainerPort{
		{
			Name:          "http",
			ContainerPort: v1alpha1.PulsarBrokerHttpServerPort,
			Protocol:      v1.ProtocolTCP,
		},
		{
			Name:          "pulsar",
			ContainerPort: v1alpha1.PulsarBrokerPulsarServerPort,
			Protocol:      v1.ProtocolTCP,
		},
	}
}

func makeContainerEnv(c *v1alpha1.Pulsar) []v1.EnvVar {
	env := make([]v1.EnvVar, 0)
	env = append(env, v1.EnvVar{
		Name:      AdvertisedAddress,
		ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{FieldPath: "status.podIP"}},
	})
	return env
}

func makeContainerEnvFrom(c *v1alpha1.Pulsar) []v1.EnvFromSource {
	froms := make([]v1.EnvFromSource, 0)

	var configRef v1.ConfigMapEnvSource
	configRef.Name = MakeConfigMapName(c)

	froms = append(froms, v1.EnvFromSource{ConfigMapRef: &configRef})
	return froms
}

func MakeWaitBrokerReadyContainer(c *v1alpha1.Pulsar) v1.Container {
	return v1.Container{
		Name:            "wait-broker-ready",
		Image:           c.Spec.Broker.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Broker.Image.PullPolicy,
		Command:         makeContainerCommand(),
		Args:            makeWaitBrokerReadyContainerCommandArgs(c),
	}
}

func makeWaitBrokerReadyContainerCommandArgs(c *v1alpha1.Pulsar) []string {
	return []string{
		fmt.Sprintf(`set -e; brokerServiceNumber="$(nslookup -timeout=10 %s-broker-service | grep Name | wc -l)"; until [ ${brokerServiceNumber} -ge 1 ]; do
			echo "pulsar cluster test-tonglinkq isn't initialized yet ... check in 10 seconds ...";
			sleep 10;
			brokerServiceNumber="$(nslookup -timeout=10 %s-broker-service | grep Name | wc -l)";
        done;`, c.GetName(), c.GetName()),
	}
}
