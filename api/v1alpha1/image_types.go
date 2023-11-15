package v1alpha1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

// ContainerImage defines the fields needed for a Docker repository image. The
// format here matches the predominant format used in Helm charts.
type ContainerImage struct {
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Repository string `json:"repository,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Tag string `json:"tag,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

func (c *ContainerImage) SetDefault(cluster *Pulsar, component string) bool {
	changed := false
	switch component {
	case ZookeeperComponent:
		if cluster.Spec.Zookeeper.Image.Repository == "" {
			cluster.Spec.Zookeeper.Image.Repository = DefaultAllPulsarContainerRepository
			changed = true
		}
		if cluster.Spec.Zookeeper.Image.Tag == "" {
			cluster.Spec.Zookeeper.Image.Tag = DefaultAllPulsarContainerVersion
			changed = true
		}
		if cluster.Spec.Zookeeper.Image.PullPolicy == "" {
			cluster.Spec.Zookeeper.Image.PullPolicy = DefaultContainerPolicy
			changed = true
		}

	case BrokerComponent:
		if cluster.Spec.Broker.Image.Repository == "" {
			cluster.Spec.Broker.Image.Repository = DefaultAllPulsarContainerRepository
			changed = true
		}
		if cluster.Spec.Broker.Image.Tag == "" {
			cluster.Spec.Broker.Image.Tag = DefaultAllPulsarContainerVersion
			changed = true
		}
		if cluster.Spec.Broker.Image.PullPolicy == "" {
			cluster.Spec.Broker.Image.PullPolicy = DefaultContainerPolicy
			changed = true
		}

	case BookieComponent:
		if cluster.Spec.Bookie.Image.Repository == "" {
			cluster.Spec.Bookie.Image.Repository = DefaultAllPulsarContainerRepository
			changed = true
		}
		if cluster.Spec.Bookie.Image.Tag == "" {
			cluster.Spec.Bookie.Image.Tag = DefaultAllPulsarContainerVersion
			changed = true
		}
		if cluster.Spec.Bookie.Image.PullPolicy == "" {
			cluster.Spec.Bookie.Image.PullPolicy = DefaultContainerPolicy
			changed = true
		}

	case AutoRecoveryComponent:
		if cluster.Spec.AutoRecovery.Image.Repository == "" {
			cluster.Spec.AutoRecovery.Image.Repository = DefaultAllPulsarContainerRepository
			changed = true
		}
		if cluster.Spec.AutoRecovery.Image.Tag == "" {
			cluster.Spec.AutoRecovery.Image.Tag = DefaultAllPulsarContainerVersion
			changed = true
		}
		if cluster.Spec.AutoRecovery.Image.PullPolicy == "" {
			cluster.Spec.AutoRecovery.Image.PullPolicy = DefaultContainerPolicy
			changed = true
		}

	//case ProxyComponent:
	//	if cluster.Spec.Proxy.Image.Repository == "" {
	//		cluster.Spec.Proxy.Image.Repository = DefaultAllPulsarContainerRepository
	//		changed = true
	//	}
	//	if cluster.Spec.Proxy.Image.Tag == "" {
	//		cluster.Spec.Proxy.Image.Tag = DefaultAllPulsarContainerVersion
	//		changed = true
	//	}
	//	if cluster.Spec.Proxy.Image.PullPolicy == "" {
	//		cluster.Spec.Proxy.Image.PullPolicy = DefaultContainerPolicy
	//		changed = true
	//	}
	//
	case ManagerComponent:
		if cluster.Spec.Manager.Image.Repository == "" {
			cluster.Spec.Manager.Image.Repository = DefaultPulsarManagerContainerRepository
			changed = true
		}
		if cluster.Spec.Manager.Image.Tag == "" {
			cluster.Spec.Manager.Image.Tag = DefaultPulsarManagerContainerVersion
			changed = true
		}
		if cluster.Spec.Manager.Image.PullPolicy == "" {
			cluster.Spec.Manager.Image.PullPolicy = corev1.PullIfNotPresent
			changed = true
		}
	case MonitorPrometheusComponent:
		if cluster.Spec.Monitor.Prometheus.Image.Repository == "" {
			cluster.Spec.Monitor.Prometheus.Image.Repository = DefaultPrometheusContainerRepository
			changed = true
		}
		if cluster.Spec.Monitor.Prometheus.Image.Tag == "" {
			cluster.Spec.Monitor.Prometheus.Image.Tag = DefaultPrometheusContainerVersion
			changed = true
		}
		if cluster.Spec.Monitor.Prometheus.Image.PullPolicy == "" {
			cluster.Spec.Monitor.Prometheus.Image.PullPolicy = corev1.PullIfNotPresent
			changed = true
		}
	case MonitorGrafanaComponent:
		if cluster.Spec.Monitor.Grafana.Image.Repository == "" {
			cluster.Spec.Monitor.Grafana.Image.Repository = DefaultMonitorGrafanaContainerRepository
			changed = true
		}
		if cluster.Spec.Monitor.Grafana.Image.Tag == "" {
			cluster.Spec.Monitor.Grafana.Image.Tag = DefaultMonitorGrafanaContainerTag
			changed = true
		}
		if cluster.Spec.Monitor.Grafana.Image.PullPolicy == "" {
			cluster.Spec.Monitor.Grafana.Image.PullPolicy = corev1.PullIfNotPresent
			changed = true
		}
	}

	return changed
}

func (c *ContainerImage) GenerateImage() string {
	return fmt.Sprintf("%s:%s", c.Repository, c.Tag)
}
