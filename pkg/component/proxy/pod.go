package proxy

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func makePodSpec(c *v1alpha1.Pulsar) v1.PodSpec {
	return v1.PodSpec{
		Affinity: c.Spec.Broker.Pod.Affinity,
		InitContainers: []v1.Container{
			makeWaitZookeeperReadyContainer(c),
			makeWaitBrokerReadyContainer(c),
		},
		Containers: []v1.Container{makeContainer(c)},
		Volumes:    makeVolumes(c),
	}
}

func makeWaitBrokerReadyContainer(c *v1alpha1.Pulsar) v1.Container {
	return v1.Container{
		Name:            "wait-broker-ready",
		Image:           c.Spec.Proxy.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Proxy.Image.PullPolicy,
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

func makeWaitZookeeperReadyContainer(c *v1alpha1.Pulsar) v1.Container {
	return v1.Container{
		Name:            "wait-zookeeper-ready",
		Image:           c.Spec.Proxy.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Proxy.Image.PullPolicy,
		Command:         makeContainerCommand(),
		Args:            makeWaitZookeeperReadyContainerCommandArgs(c),
		EnvFrom:         makeWaitZookeeperReadyContainerEnvFrom(c),
	}
}

func makeWaitZookeeperReadyContainerCommandArgs(c *v1alpha1.Pulsar) []string {
	return []string{
		fmt.Sprintf(`until bin/pulsar zookeeper-shell -server %s-zookeeper-service get /admin/clusters/%s; do
		sleep 3;
		done;`, c.GetName(), c.GetName()),
	}
}

func makeWaitZookeeperReadyContainerEnvFrom(c *v1alpha1.Pulsar) []v1.EnvFromSource {
	return []v1.EnvFromSource{
		{ConfigMapRef: &v1.ConfigMapEnvSource{LocalObjectReference: v1.LocalObjectReference{Name: MakeConfigMapName(c)}}},
	}
}

func makeContainer(c *v1alpha1.Pulsar) v1.Container {
	return v1.Container{
		Name:            "proxy",
		Image:           c.Spec.Proxy.Image.GenerateImage(),
		ImagePullPolicy: c.Spec.Proxy.Image.PullPolicy,
		Command:         makeContainerCommand(),
		Args:            makeContainerCommandArgs(),
		Ports:           makeContainerPort(c),
		EnvFrom:         makeContainerEnvFrom(c),
		ReadinessProbe: &v1.Probe{
			InitialDelaySeconds: 30,
			PeriodSeconds:       10,
			FailureThreshold:    10,
			SuccessThreshold:    1,
			TimeoutSeconds:      1,
			ProbeHandler: v1.ProbeHandler{
				HTTPGet: &v1.HTTPGetAction{
					Path: "/status.html",
					Port: intstr.IntOrString{
						IntVal: v1alpha1.PulsarBrokerHttpServerPort,
					},
					Scheme: v1.URISchemeHTTP,
				},
			},
		},
		LivenessProbe: &v1.Probe{
			InitialDelaySeconds: 30,
			PeriodSeconds:       10,
			FailureThreshold:    10,
			SuccessThreshold:    1,
			TimeoutSeconds:      1,
			ProbeHandler: v1.ProbeHandler{
				HTTPGet: &v1.HTTPGetAction{
					Path: "/status.html",
					Port: intstr.IntOrString{
						IntVal: v1alpha1.PulsarBrokerHttpServerPort,
					},
					Scheme: v1.URISchemeHTTP,
				},
			},
		},
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      makeLog4j2Name(c),
				MountPath: "/pulsar/conf/log4j2.yaml",
				SubPath:   "log4j2.yaml",
			},
			{
				Name:      "token-keys",
				MountPath: "/pulsar/keys",
				ReadOnly:  true,
			},
			{
				Name:      "proxy-token",
				MountPath: "/pulsar/tokens",
				ReadOnly:  true,
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
		"bin/apply-config-from-env.py conf/proxy.conf; echo 'OK' > status; bin/pulsar proxy;",
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

func makeContainerEnvFrom(c *v1alpha1.Pulsar) []v1.EnvFromSource {
	froms := make([]v1.EnvFromSource, 0)

	var configRef v1.ConfigMapEnvSource
	configRef.Name = MakeConfigMapName(c)

	froms = append(froms, v1.EnvFromSource{ConfigMapRef: &configRef})
	return froms
}

func makeLog4j2Name(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-proxy-log4j2", c.Namespace)
}

func makeVolumes(c *v1alpha1.Pulsar) []v1.Volume {
	var defaultMode int32 = 420
	return []v1.Volume{
		{
			Name: makeLog4j2Name(c),
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{Name: MakeConfigMapName(c)},
					DefaultMode:          &defaultMode,
				},
			},
		},
		{
			Name: "token-keys",
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: fmt.Sprintf("%s-token-asymmetric-key", c.Name),
					Items: []v1.KeyToPath{
						{
							Key:  "PUBLICKEY",
							Path: "token/public.key",
						},
					},
					DefaultMode: &defaultMode,
				},
			},
		},
		{
			Name: "proxy-token",
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: fmt.Sprintf("%s-token-proxy-admin", c.Name),
					Items: []v1.KeyToPath{
						{
							Key:  "TOKEN",
							Path: "proxy/token",
						},
					},
					DefaultMode: &defaultMode,
				},
			},
		},
	}
}
