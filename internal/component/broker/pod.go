package broker

import (
	"fmt"
	"strconv"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/internal/component/bookie"
	"github.com/falser101/pulsar-operator/internal/component/zookeeper"
	v1 "k8s.io/api/core/v1"
)

func makePodSpec(c *v1alpha1.PulsarCluster) v1.PodSpec {
	var podSepc = v1.PodSpec{
		Affinity: c.Spec.Broker.Pod.Affinity,
		InitContainers: []v1.Container{
			zookeeper.MakeWaitZookeeperReadyContainer(c),
			makeWaitBookieReadyContainer(c),
		},
		Containers: []v1.Container{makeContainer(c)},
	}
	if c.Spec.Auth.AuthenticationEnabled {
		podSepc.Volumes = MakeAuthenticationVolumes(c)
	}
	return podSepc
}

func MakeAuthenticationVolumes(c *v1alpha1.PulsarCluster) []v1.Volume {
	return []v1.Volume{
		{
			Name: "token-keys",
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: MakeAsymmetricKey(c),
					Items: []v1.KeyToPath{
						{
							Key:  "PUBLICKEY",
							Path: "token/public.key",
						},
					},
				},
			},
		},
		{
			Name: "broker-token",
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: MakeBrokerTokenName(c),
					Items: []v1.KeyToPath{
						{
							Key:  "TOKEN",
							Path: "broker/token",
						},
					},
				},
			},
		},
	}
}

func MakeBrokerTokenName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-token-broker-admin", c.Name)
}

func MakeAsymmetricKey(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-token-asymmetric-key", c.Name)
}

func makeWaitBookieReadyContainer(c *v1alpha1.PulsarCluster) v1.Container {
	return v1.Container{
		Name:            "wait-bookie-ready",
		Image:           c.Spec.Broker.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Broker.Image.PullPolicy,
		Command:         []string{"sh", "-c"},
		EnvFrom: []v1.EnvFromSource{
			{
				ConfigMapRef: &v1.ConfigMapEnvSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: bookie.MakeConfigMapName(c),
					},
				},
			},
		},
		Args: []string{
			fmt.Sprintf(`export BOOKIE_MEM="-Xmx128M";
			bin/apply-config-from-env.py conf/bookkeeper.conf;
            until bin/bookkeeper shell whatisinstanceid; do
              echo "bookkeeper cluster is not initialized yet. backoff for 3 seconds ...";
              sleep 3;
            done;
            echo "bookkeeper cluster is already initialized";
            bookieServiceNumber="$(nslookup -timeout=10 %s | grep Name | wc -l)";
            until [ ${bookieServiceNumber} -ge %s ]; do
              echo "bookkeeper cluster %s isn't ready yet ... check in 10 seconds ...";
              sleep 10;
              bookieServiceNumber="$(nslookup -timeout=10 %s | grep Name | wc -l)";
            done;
            echo "bookkeeper cluster is ready";`, bookie.MakeServiceName(c), c.Spec.Broker.ConfigData["managedLedgerDefaultEnsembleSize"], c.Name, bookie.MakeServiceName(c)),
		},
	}
}

func makeContainer(c *v1alpha1.PulsarCluster) v1.Container {
	var container = v1.Container{
		Name:            "broker",
		Image:           c.Spec.Broker.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Broker.Image.PullPolicy,
		Command:         makeContainerCommand(),
		Args:            makeContainerCommandArgs(c),
		Ports:           makeContainerPort(c),
		Env:             makeContainerEnv(c),
		EnvFrom:         makeContainerEnvFrom(c),
	}
	if c.Spec.Auth.AuthenticationEnabled {
		container.VolumeMounts = append(
			container.VolumeMounts,
			v1.VolumeMount{
				Name:      "token-keys",
				ReadOnly:  true,
				MountPath: "/pulsar/keys",
			},
			v1.VolumeMount{
				Name:      "broker-token",
				ReadOnly:  true,
				MountPath: "/pulsar/tokens",
			},
		)
	}
	return container
}

func makeContainerCommand() []string {
	return []string{
		"sh",
		"-c",
	}
}

func makeContainerCommandArgs(c *v1alpha1.PulsarCluster) []string {
	return []string{
		fmt.Sprintf(`bin/apply-config-from-env.py conf/broker.conf;
		bin/gen-yml-from-env.py conf/functions_worker.yml;
		echo "OK" > "${statusFilePath:-status}";
		bin/pulsar zookeeper-shell -server %s get %s;
		while [ $? -eq 0 ]; do
		  echo "broker %s znode still exists ... check in 10 seconds ...";
		  sleep 10;
		  bin/pulsar zookeeper-shell -server %s get %s;
		done;
		cat conf/pulsar_env.sh;
		OPTS="${OPTS} -Dlog4j2.formatMsgNoLookups=true" exec bin/pulsar broker;`, zookeeper.Connect(c), znode(c), hostname(c), zookeeper.Connect(c), znode(c)),
	}
}

func makeContainerPort(c *v1alpha1.PulsarCluster) []v1.ContainerPort {
	http, _ := strconv.Atoi(c.Spec.Broker.Ports.Http)
	pulsar, _ := strconv.Atoi(c.Spec.Broker.Ports.Pulsar)
	return []v1.ContainerPort{
		{
			Name:          "http",
			ContainerPort: int32(http),
			Protocol:      v1.ProtocolTCP,
		},
		{
			Name:          "pulsar",
			ContainerPort: int32(pulsar),
			Protocol:      v1.ProtocolTCP,
		},
	}
}

func makeContainerEnv(c *v1alpha1.PulsarCluster) []v1.EnvVar {
	return []v1.EnvVar{
		{
			Name:      AdvertisedAddress,
			ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{FieldPath: "status.podIP"}},
		}}
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

func MakeWaitBrokerReadyContainer(c *v1alpha1.PulsarCluster) v1.Container {
	return v1.Container{
		Name:            "wait-broker-ready",
		Image:           c.Spec.Broker.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Broker.Image.PullPolicy,
		Command:         makeContainerCommand(),
		Args:            makeWaitBrokerReadyContainerCommandArgs(c),
	}
}

func makeWaitBrokerReadyContainerCommandArgs(c *v1alpha1.PulsarCluster) []string {
	return []string{
		"set -e;",
		fmt.Sprintf("brokerServiceNumber=\"$(nslookup -timeout=10 %s-broker-service | grep Name | wc -l)\";", c.Name),
		"until [ ${brokerServiceNumber} -ge 1 ];",
		"do echo \"pulsar cluster test-tonglinkq isn't initialized yet ... check in 10 seconds ...\";",
		"sleep 10;",
		fmt.Sprintf("brokerServiceNumber=\"$(nslookup -timeout=10 %s-broker-service | grep Name | wc -l)\";", c.Name),
		"done;",
	}
}

func znode(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("/loadbalance/brokers/%s:%s", hostname(c), c.Spec.Broker.Ports.Http)
}
